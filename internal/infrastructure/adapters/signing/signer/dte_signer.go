package signer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type DTESigner struct {
	clientRepo auth.AuthRepositoryPort
	client     *http.Client
}

type SignRequest struct {
	NIT         string          `json:"nit"`
	Activo      bool            `json:"activo"`
	PasswordPri string          `json:"passwordPri"`
	DteJson     json.RawMessage `json:"dteJson"`
}

type SignResponse struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// SpringBootError representa el formato estándar de error de Spring Boot
type SpringBootError struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

func NewDTESigner(clientRepo auth.AuthRepositoryPort) *DTESigner {
	return &DTESigner{
		clientRepo: clientRepo,
		client:     &http.Client{Timeout: 2 * time.Second},
	}
}

func (s *DTESigner) SignDTE(ctx context.Context, dte json.RawMessage, nit string) (string, error) {
	client, err := s.clientRepo.GetByNIT(ctx, nit)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("DTESigner", "SignDTE", "Error getting client by NIT", err)
	}

	req := SignRequest{
		NIT:         nit,
		Activo:      true,
		PasswordPri: client.PasswordPri,
		DteJson:     dte,
	}

	logs.Debug("Signing request", map[string]interface{}{
		"nit":     nit,
		"DteJson": string(dte),
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshalling sign request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		"POST",
		config.Signer.Path,
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating sign request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error calling signer service: %w", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var springError SpringBootError
		if err := json.Unmarshal(body, &springError); err == nil {
			// Si podemos deserializar como error de Spring Boot, usamos ese formato
			return "", fmt.Errorf("signer service error - status: %d, error: %s, message: %s",
				springError.Status,
				springError.Error,
				springError.Message)
		}

		// Si no es un error de Spring Boot, intentamos el formato de error del firmador
		var errorResp struct {
			Status string `json:"status"`
			Body   struct {
				Codigo  string   `json:"codigo"`
				Mensaje []string `json:"mensaje"`
			} `json:"body"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return "", fmt.Errorf("signer service error - code: %s, messages: %v",
				errorResp.Body.Codigo,
				errorResp.Body.Mensaje)
		}

		// Si ninguno de los formatos coincide, devolvemos el body como string
		return "", fmt.Errorf("unexpected error response from signer service (status %d): %s",
			resp.StatusCode, string(body))
	}

	// Para respuestas exitosas
	var signResp SignResponse
	if err := json.Unmarshal(body, &signResp); err != nil {
		return "", fmt.Errorf("error decoding successful response: %w, body: %s", err, string(body))
	}

	logs.Debug("Sign response", map[string]interface{}{
		"status": signResp.Body,
	})
	logs.Info("Document signed successfully")

	return signResp.Body, nil
}
