package batch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config/env"
	authPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	authModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/circuit"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/hacienda_error"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// BatchTransmitterService implementa la lógica de transmisión de lotes a Hacienda
type BatchTransmitterService struct {
	haciendaAuth   authPorts.HaciendaAuthManager
	signer         authPorts.SignerManager
	dteRepo        ports.DTERepositoryPort
	httpClient     *http.Client
	config         models.BatchConfig
	timeProvider   interfaces.TimeProvider
	circuitBreaker *circuit.CircuitBreaker
}

// NewBatchTransmitterService constructor para BatchTransmitterService
func NewBatchTransmitterService(
	haciendaAuth authPorts.HaciendaAuthManager,
	signer authPorts.SignerManager,
	dteRepo ports.DTERepositoryPort,
	config models.BatchConfig,
	timeProvider interfaces.TimeProvider,
) *BatchTransmitterService {
	return &BatchTransmitterService{
		haciendaAuth: haciendaAuth,
		signer:       signer,
		dteRepo:      dteRepo,
		config:       config,
		timeProvider: timeProvider,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		},
		circuitBreaker: circuit.NewCircuitBreaker(
			3,
			5*time.Minute,
		),
	}
}

// GetDTEVersion determina la versión según el tipo de DTE
func (s *BatchTransmitterService) GetDTEVersion(dteType string) int {
	switch dteType {
	case constants.FacturaElectronica:
		return 1
	case constants.CCFElectronico:
		return 2
	default:
		return 1 // Versión por defecto
	}
}

// TransmitBatch transmite un lote de documentos a Hacienda
func (s *BatchTransmitterService) TransmitBatch(
	ctx context.Context,
	systemNIT string,
	dteType string,
	signedDocs []string,
	token string,
	creds authModels.HaciendaCredentials,
) (*models.BatchResponse, error) {
	if len(signedDocs) == 0 {
		return nil, errors.New("no documents to transmit")
	}

	batchID := strings.ToUpper(uuid.New().String())
	batch := &models.BatchRequest{
		Ambient:   s.config.GetAmbient(),
		SendID:    batchID,
		Version:   s.GetDTEVersion(dteType),
		NIT:       systemNIT,
		Documents: signedDocs,
	}

	haciendaToken, err := s.getHaciendaTokenWithRetry(ctx, token, creds)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errPackage.ErrHaciendaTokenGeneration, err)
	}

	response, err := s.sendBatchWithRetry(ctx, batch, haciendaToken)
	if err != nil {
		logs.Error("Failed to send batch", map[string]interface{}{
			"error":   err.Error(),
			"batchId": batchID,
		})
		return nil, err
	}

	logs.Info("Batch sent successfully", map[string]interface{}{
		"batchId": response.BatchCode,
		"status":  response.Status,
		"msg":     response.Description,
		"idEnvio": response.SendID,
	})

	return response, nil
}

// getHaciendaTokenWithRetry obtiene un token de autenticación de Hacienda con reintentos
func (s *BatchTransmitterService) getHaciendaTokenWithRetry(
	ctx context.Context,
	token string,
	creds authModels.HaciendaCredentials,
) (string, error) {
	var haciendaToken string
	var err error
	retryPolicy := s.config.GetRetryPolicy()

	for attempt := 0; attempt < retryPolicy.MaxAttempts; attempt++ {
		haciendaToken, err = s.haciendaAuth.GetOrCreateHaciendaTokenWithCreds(ctx, token, creds)
		if err == nil {
			return haciendaToken, nil
		}

		if !s.shouldRetry(err) {
			return "", err
		}

		s.sleep(attempt)
	}

	return "", fmt.Errorf("max retry attempts reached: %w", err)
}

// sendBatchWithRetry envía un lote con reintentos
func (s *BatchTransmitterService) sendBatchWithRetry(
	ctx context.Context,
	batch *models.BatchRequest,
	token string,
) (*models.BatchResponse, error) {
	var response *models.BatchResponse
	var err error
	retryPolicy := s.config.GetRetryPolicy()

	for attempt := 0; attempt < retryPolicy.MaxAttempts; attempt++ {
		response, err = s.transmitToHacienda(ctx, batch, token)
		if err == nil {
			return response, nil
		}

		if !s.shouldRetry(err) {
			return nil, err
		}

		s.sleep(attempt)
	}

	return nil, fmt.Errorf("max retry attempts reached: %w", err)
}

// transmitToHacienda envía el lote a Hacienda
func (s *BatchTransmitterService) transmitToHacienda(
	ctx context.Context,
	batch *models.BatchRequest,
	token string,
) (*models.BatchResponse, error) {
	if !s.circuitBreaker.AllowRequest() {
		logs.Warn("Circuit breaker preventing request to Hacienda", map[string]interface{}{
			"state": s.circuitBreaker.GetState(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"BatchTransmitterService",
			"transmitToHacienda",
			"service temporarily unavailable due to consecutive failures",
			nil,
		)
	}

	logs.Info("Sending batch to Hacienda", map[string]interface{}{
		"batchId": batch.SendID,
		"ambient": batch.Ambient,
		"docs":    len(batch.Documents),
	})

	reqBody, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", env.MHPaths.LoteReceptionURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.circuitBreaker.RecordFailure()
		logs.Error("Request to Hacienda failed", map[string]interface{}{
			"error":        err.Error(),
			"failureCount": s.circuitBreaker.GetFailureCount(),
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("hacienda returned status: %d", resp.StatusCode)
	}

	var batchResp models.BatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		s.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.circuitBreaker.RecordSuccess()
	return &batchResp, nil
}

// VerifyBatchStatus verifica el estado de un lote
func (s *BatchTransmitterService) VerifyBatchStatus(
	ctx context.Context,
	batchID string,
	mhBatchID string,
	token string,
	docsMap map[string]dte.DTEDetails,
) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	deadline := utils.TimeNow().Add(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			status, isProcessed, err := s.checkBatchStatus(ctx, mhBatchID, token)
			if err != nil {
				logs.Error("Failed to check batch status, verify", map[string]interface{}{
					"error":   err.Error(),
					"batchID": mhBatchID,
				})
				return fmt.Errorf("failed to check batch status err: %w", err)
			}

			if !isProcessed {
				logs.Info("Batch still processing", map[string]interface{}{
					"batchID": mhBatchID,
				})

				if utils.TimeNow().After(deadline) {
					return fmt.Errorf("timeout waiting for batch processing")
				}
				continue
			}

			// Procesar documentos procesados
			if len(status.Processed) > 0 {
				for _, processed := range status.Processed {
					if doc, exists := docsMap[processed.GenerationCode]; exists {
						logs.Info("Document processed", map[string]interface{}{
							"code":            processed.MessageCode,
							"message":         processed.DescriptionMessage,
							"observations":    processed.Observations,
							"processedAt":     processed.ProcessingDate,
							"reception_stamp": processed.ReceptionStamp,
						})

						if err := s.registerProcessedDocument(ctx, doc, constants.DocumentReceived, processed.ReceptionStamp); err != nil {
							logs.Error("Failed to register received document in dte_documents", map[string]interface{}{
								"error": err.Error(),
								"id":    doc.ID,
							})
						}
					}
				}
			}

			// Procesar documentos rechazados
			if len(status.Rejected) > 0 {
				for _, rejected := range status.Rejected {
					if doc, exists := docsMap[rejected.GenerationCode]; exists {
						logs.Info("Document rejected", map[string]interface{}{
							"code":         rejected.MessageCode,
							"message":      rejected.DescriptionMessage,
							"observations": rejected.Observations,
							"processedAt":  rejected.ProcessingDate,
						})

						if err := s.registerProcessedDocument(ctx, doc, constants.DocumentRejected, ""); err != nil {
							logs.Error("Failed to register document in dte_documents", map[string]interface{}{
								"error": err.Error(),
								"id":    doc.ID,
							})
						}
					}
				}
			}

			logs.Info("Batch status verified", map[string]interface{}{
				"batchID":        batchID,
				"totalProcessed": len(status.Processed),
				"totalRejected":  len(status.Rejected),
			})

			return nil
		}
	}
}

// checkBatchStatus verifica el estado de un lote en Hacienda
func (s *BatchTransmitterService) checkBatchStatus(ctx context.Context, batchID string, haciendaToken string) (*models.ConsultBatchResponse, bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/%s", env.MHPaths.LoteReceptionConsultURL, batchID),
		nil,
	)
	if err != nil {
		logs.Error("Failed to create batch status request", map[string]interface{}{
			"error":   err.Error(),
			"batchID": batchID,
		})
		return nil, false, err
	}

	req.Header.Set("Authorization", haciendaToken)
	req.Header.Set("Content-Type", "application/json")

	logs.Info("Checking batch status", map[string]interface{}{
		"url":     req.URL.String(),
		"method":  req.Method,
		"batchID": batchID,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error("Failed to check batch status inner", map[string]interface{}{
			"error":   err.Error(),
			"batchID": batchID,
		})
		return nil, false, err
	}
	defer resp.Body.Close()

	var batchResp models.ConsultBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		if err == io.EOF {
			// Si no hay contenido, el lote aun no ha sido procesado
			return nil, false, nil
		}
		logs.Error("Failed to decode batch response", map[string]interface{}{
			"error":   err.Error(),
			"batchID": batchID,
		})
		return nil, false, err
	}

	logs.Info("Batch status response", map[string]interface{}{
		"Processed": len(batchResp.Processed),
		"Rejected":  len(batchResp.Rejected),
	})

	return &batchResp, true, nil
}

// registerProcessedDocument registra un documento procesado en el repositorio
func (s *BatchTransmitterService) registerProcessedDocument(ctx context.Context, doc dte.DTEDetails, status string, stamp string) error {
	err := s.dteRepo.Update(ctx, doc.ID, status, utils.ToStringPointer(stamp))
	if err != nil {
		logs.Error("Failed to register document in dte_documents", map[string]interface{}{
			"error": err.Error(),
			"id":    doc.ID,
		})
		return err
	}

	return nil
}

// shouldRetry determina si se debe reintentar una operación
func (s *BatchTransmitterService) shouldRetry(err error) bool {
	// Errores de red/conexión - siempre reintentar
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		logs.Info("Network error detected, will retry", map[string]interface{}{
			"error": err.Error(),
		})
		return true
	}

	var httpErr *hacienda_error.HTTPResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode >= 500 && httpErr.StatusCode <= 599 {
			logs.Info("Server error detected, will retry", map[string]interface{}{
				"statusCode": httpErr.StatusCode,
			})
			return true
		}

		switch httpErr.StatusCode {
		case http.StatusTooManyRequests, // 429
			http.StatusRequestTimeout,     // 408
			http.StatusBadGateway,         // 502
			http.StatusServiceUnavailable, // 503
			http.StatusGatewayTimeout:     // 504
			logs.Info("Retryable HTTP error detected", map[string]interface{}{
				"statusCode": httpErr.StatusCode,
			})
			return true
		}

		// No reintentar otros codigos HTTP
		logs.Info("Non-retryable HTTP error", map[string]interface{}{
			"statusCode": httpErr.StatusCode,
		})
		return false
	}

	var haciendaErr *hacienda_error.HaciendaResponseError
	if errors.As(err, &haciendaErr) {
		// No reintentar errores de validación o autorización
		if strings.Contains(strings.ToLower(haciendaErr.Description), "validaci") ||
			strings.Contains(strings.ToLower(haciendaErr.Description), "autoriza") {
			logs.Info("Non-retryable Hacienda error", map[string]interface{}{
				"code":    haciendaErr.Code,
				"message": haciendaErr.Description,
			})
			return false
		}

		// Reintentar otros errores de Hacienda
		logs.Info("Retryable Hacienda error", map[string]interface{}{
			"code":    haciendaErr.Code,
			"message": haciendaErr.Description,
		})
		return true
	}

	// Errores de contexto
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		logs.Info("Context error detected, will retry", map[string]interface{}{
			"error": err.Error(),
		})
		return true
	}

	// Para errores no clasificados, no son retryables
	// pero se loguea para análisis
	logs.Warn("Unclassified error, defaulting to retry", map[string]interface{}{
		"error": err.Error(),
	})
	return true
}

// sleep implementa el backoff exponencial para reintentos
func (s *BatchTransmitterService) sleep(attempt int) {
	retryPolicy := s.config.GetRetryPolicy()
	backoff := retryPolicy.InitialInterval * time.Duration(float64(attempt)*retryPolicy.BackoffFactor)
	if backoff > retryPolicy.MaxInterval {
		backoff = retryPolicy.MaxInterval
	}
	s.timeProvider.Sleep(backoff)
}
