package transmitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/hacienda_error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/processors"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
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
	HaciendaToken string
	haciendaAuth  ports.HaciendaAuthManager
	httpClient    *http.Client
	processors    map[string]DocumentProcessor
}

func NewMHTransmitter(haciendaAuth ports.HaciendaAuthManager) ports.DTETransmitter {
	t := &MHTransmitter{
		haciendaAuth: haciendaAuth,
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

func (t *MHTransmitter) Transmit(ctx context.Context, document interface{}, signedDoc string, systemToken string) (*models.TransmitResult, error) {
	// Forzar modo de contingencia si está activado
	if env.Server.ForceContingency && env.Server.AmbientCode == "00" {
		logs.Info("Forcing contingency mode - simulating service unavailable")
		return nil, &hacienda_error.HTTPResponseError{
			StatusCode: http.StatusServiceUnavailable,
			Body:       []byte("Forced contingency - service unavailable"),
			URL:        env.MHPaths.ReceptionURL,
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
		return nil, err
	}

	// Procesar respuesta
	return processor.ProcessResponse(resp)
}

func (t *MHTransmitter) CheckDocumentStatus(ctx context.Context, document interface{}, nit string) (*models.TransmitResult, error) {
	_, dteType, generationCode, _, err := processors.GetDocumentRequestData(document)

	haciendaReqBody := HaciendaConsultRequest{
		IssuerNIT:      nit,
		DTEType:        dteType,
		GenerationCode: generationCode,
	}

	jsonData, err := json.Marshal(haciendaReqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", env.MHPaths.ReceptionConsultURL, bytes.NewBuffer(jsonData))
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

	var haciendaResp models.HaciendaResponse
	if err := json.NewDecoder(resp.Body).Decode(&haciendaResp); err != nil {
		return nil, err
	}

	return &models.TransmitResult{
		Status:         haciendaResp.Status,
		ReceptionStamp: &haciendaResp.ReceptionStamp,
		ProcessingDate: haciendaResp.ProcessingDate,
		MessageCode:    haciendaResp.MessageCode,
		MessageDesc:    haciendaResp.DescriptionMessage,
		Observations:   haciendaResp.Observations,
	}, nil
}

func (t *MHTransmitter) SendToHacienda(ctx context.Context, request *models.HaciendaRequest, systemToken string) (*models.HaciendaResponse, error) {
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
		url = env.MHPaths.ReceptionURL
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

	var haciendaResp models.HaciendaResponse
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

	var response models.HaciendaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w, body: %s", err, string(body))
	}

	return &response, nil
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
