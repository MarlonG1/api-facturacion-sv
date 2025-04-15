package transmitter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	models2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	ports2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/hacienda_error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/processors"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type HaciendaConsultRequest struct {
	IssuerNIT      string `json:"nitEmisor"`
	DTEType        string `json:"tdte"`
	GenerationCode string `json:"codigoGeneracion"`
}

type MHTransmitter struct {
	HaciendaToken      string
	haciendaAuth       ports.HaciendaAuthManager
	failedSequenceRepo ports2.FailedSequenceNumberRepositoryPort
	httpClient         *http.Client
	processors         map[string]DocumentProcessor
}

func NewMHTransmitter(haciendaAuth ports.HaciendaAuthManager, failedSequenceRepo ports2.FailedSequenceNumberRepositoryPort) ports.DTETransmitter {
	t := &MHTransmitter{
		haciendaAuth:       haciendaAuth,
		failedSequenceRepo: failedSequenceRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		processors: make(map[string]DocumentProcessor),
	}

	// Registrar processors
	t.processors["invalidation"] = &processors.InvalidationProcessor{}
	t.processors["dte"] = &processors.DTEProcessor{}
	return t
}

func (t *MHTransmitter) Transmit(ctx context.Context, document interface{}, signedDoc string, systemToken string) (*models2.TransmitResult, error) {
	// Forzar modo de contingencia si está activado
	if config.Server.ForceContingency && config.Server.AmbientCode == "00" {
		logs.Info("Forcing contingency mode - simulating service unavailable")
		return nil, &hacienda_error.HTTPResponseError{
			StatusCode: http.StatusServiceUnavailable,
			Body:       []byte("Forced contingency - service unavailable"),
			URL:        config.MHPaths.ReceptionURL,
			Method:     "POST",
		}
	}

	processor := t.getProcessor(document)
	if processor == nil {
		return nil, fmt.Errorf("no processor found for document type: %T", document)
	}

	// Preparar request
	req, err := processor.ProcessRequest(signedDoc, document)
	if err != nil {
		logs.Error("Failed to process request", map[string]interface{}{
			"error": err.Error(),
			"type":  fmt.Sprintf("%T", signedDoc),
		})
		return nil, err
	}

	// Enviar a Hacienda
	resp, err := t.SendToHacienda(ctx, req, systemToken)
	if err != nil {
		t.handleFailedSequence(ctx, err, document, req)
		return nil, err
	}

	// Procesar respuesta
	return processor.ProcessResponse(resp)
}

func (t *MHTransmitter) CheckDocumentStatus(ctx context.Context, document interface{}, nit string) (*models2.TransmitResult, error) {
	_, dteType, generationCode, _, err := processors.GetDocumentRequestData(document)

	haciendaReqBody := HaciendaConsultRequest{
		IssuerNIT:      nit,
		DTEType:        dteType,
		GenerationCode: generationCode,
	}

	jsonData, err := json.Marshal(haciendaReqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", config.MHPaths.ReceptionConsultURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logs.Info("Failed to create request", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	req.Header.Set("Authorization", t.HaciendaToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logs.Error("Failed to check document status", map[string]interface{}{
			"error": err.Error(),
			"code":  generationCode,
		})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		logs.Error("Failed to check document status", map[string]interface{}{
			"status": resp.Status,
		})
		return nil, fmt.Errorf("failed to check document status: %s", resp.Status)
	}

	var haciendaResp models2.HaciendaResponse
	if err := json.NewDecoder(resp.Body).Decode(&haciendaResp); err != nil {
		return nil, err
	}

	return &models2.TransmitResult{
		Status:         haciendaResp.Status,
		ReceptionStamp: &haciendaResp.ReceptionStamp,
		ProcessingDate: haciendaResp.ProcessingDate,
		MessageCode:    haciendaResp.MessageCode,
		MessageDesc:    haciendaResp.DescriptionMessage,
		Observations:   haciendaResp.Observations,
	}, nil
}

func (t *MHTransmitter) SendToHacienda(ctx context.Context, request *models2.HaciendaRequest, systemToken string) (*models2.HaciendaResponse, error) {
	err := t.getHaciendaToken(ctx, systemToken)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := request.URL
	if url == "" {
		url = config.MHPaths.ReceptionURL
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", t.HaciendaToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "HaciendaApp/1.0")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logs.Error("Failed to send to Hacienda", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var haciendaResp models2.HaciendaResponse
	if err := json.Unmarshal(body, &haciendaResp); err == nil {

		if haciendaResp.Status == "RECHAZADO" {
			logs.Error("Document rejected by Hacienda", map[string]interface{}{
				"code":         haciendaResp.MessageCode,
				"message":      haciendaResp.DescriptionMessage,
				"status":       haciendaResp.Status,
				"observations": haciendaResp.Observations,
				"processedAt":  haciendaResp.ProcessingDate,
			})
			return nil, hacienda_error.NewHaciendaError(&haciendaResp, resp.StatusCode)
		}
		return &haciendaResp, nil
	}

	var response models2.HaciendaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w, body: %s", err, string(body))
	}

	return &response, nil
}

func (t *MHTransmitter) handleFailedSequence(ctx context.Context, err error, document interface{}, req *models2.HaciendaRequest) {
	// Only register if it's a Hacienda error
	var haciendaErr *hacienda_error.HaciendaResponseError
	if !errors.As(err, &haciendaErr) {
		return
	}

	// Extract necessary information for failed sequence
	claims, ok := ctx.Value("claims").(*models.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context for failed sequence", nil)
		return
	}

	// Extract sequence number from control number
	_, dteType, _, sequenceNumber, extractErr := processors.GetDocumentRequestData(document)
	if extractErr != nil {
		logs.Error("Failed to extract document data for failed sequence", map[string]interface{}{
			"error": extractErr.Error(),
		})
		return
	}

	// Get current year
	currentYear := uint(utils.TimeNow().Year())

	// Register the failed sequence
	registrationErr := t.failedSequenceRepo.RegisterFailedSequence(
		ctx,
		claims.BranchID,
		dteType,
		uint(sequenceNumber),
		currentYear,
		haciendaErr.Description,
		haciendaErr.Code,
		document,
		formatHaciendaErrorResponse(haciendaErr),
	)

	if registrationErr != nil {
		logs.Error("Failed to register failed sequence", map[string]interface{}{
			"error":          registrationErr.Error(),
			"branchID":       claims.BranchID,
			"dteType":        dteType,
			"sequenceNumber": sequenceNumber,
		})
	} else {
		logs.Info("Failed sequence registered successfully", map[string]interface{}{
			"branchID":       claims.BranchID,
			"dteType":        dteType,
			"sequenceNumber": sequenceNumber,
			"errorCode":      haciendaErr.Code,
		})
	}
}

func formatHaciendaErrorResponse(err *hacienda_error.HaciendaResponseError) string {
	// Convert the Hacienda error to a JSON string for storage
	response := map[string]interface{}{
		"status":         err.Status,
		"code":           err.Code,
		"description":    err.Description,
		"classification": err.Classification,
		"statusCode":     err.StatusCode,
		"observations":   err.Observations,
		"processedAt":    err.ProcessedAt,
	}

	jsonData, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logs.Error("Failed to marshal Hacienda error response", map[string]interface{}{
			"error": jsonErr.Error(),
		})
		return fmt.Sprintf("Error parsing response: %s", err.Error())
	}

	return string(jsonData)
}

func (t *MHTransmitter) getProcessor(document interface{}) DocumentProcessor {
	var docMap map[string]interface{}

	// Si es string JSON, parsearlo
	if jsonStr, ok := document.(string); ok {
		if err := json.Unmarshal([]byte(jsonStr), &docMap); err != nil {
			logs.Error("Failed to parse document JSON", map[string]interface{}{
				"error": err.Error(),
			})
			return nil
		}
	} else {
		// Si no es string, intentar convertir directamente
		jsonBytes, err := json.Marshal(document)
		if err != nil {
			logs.Error("Failed to marshal document", map[string]interface{}{
				"error": err.Error(),
			})
			return nil
		}
		if err := json.Unmarshal(jsonBytes, &docMap); err != nil {
			logs.Error("Failed to parse marshaled document", map[string]interface{}{
				"error": err.Error(),
			})
			return nil
		}
	}

	// Determinar el tipo de documento basado en su estructura
	if isInvalidationDocument(docMap) {
		logs.Info("Document identified as invalidation")
		return t.processors["invalidation"]
	}

	// Si no es invalidación, se asume que es un DTE
	logs.Info("Document identified as DTE")
	return t.processors["dte"]
}

func (t *MHTransmitter) getHaciendaToken(ctx context.Context, systemToken string) error {
	haciendaToken, err := t.haciendaAuth.GetOrCreateHaciendaToken(ctx, systemToken)
	if err != nil {
		logs.Error("Failed to get Hacienda token", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	t.HaciendaToken = haciendaToken
	return nil
}

func isInvalidationDocument(doc map[string]interface{}) bool {
	if motivo, exists := doc["motivo"].(map[string]interface{}); exists {
		_, hasType := motivo["tipoAnulacion"]
		_, hasResponsible := motivo["nombreResponsable"]
		return hasType && hasResponsible
	}

	_, hasSummary := doc["resumen"]
	_, hasItems := doc["cuerpoDocumento"]

	// Si tiene resumen o items, NO es una invalidación
	return !(hasSummary || hasItems)
}
