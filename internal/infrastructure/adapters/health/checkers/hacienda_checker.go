package checkers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/dimiro1/health"
	"io"
	"net/http"
	"time"
)

// haciendaChecker implementa tanto la interfaz ports.ComponentChecker como health.Checker
type haciendaChecker struct {
	client *http.Client
}

func NewHaciendaChecker() ports.ComponentChecker {
	return &haciendaChecker{
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *haciendaChecker) Name() string {
	return "hacienda"
}

// Check implementa ports.ComponentChecker.Check
func (c *haciendaChecker) Check() models.Health {
	// CustomHealthChecker para Hacienda
	health := c.checkHealth()

	status := constants.StatusUp
	details := "Hacienda service is healthy"

	if health.IsDown() {
		status = constants.StatusDown
		details = "Hacienda service is down"

		// Extraer detalles si están disponibles
		if health.GetInfo("error") != nil {
			details = fmt.Sprintf("%s: %v", details, health.GetInfo("error"))
		}
	}

	return models.Health{
		Status:  status,
		Details: details,
	}
}

func (c *haciendaChecker) checkHealth() health.Health {
	result := health.NewHealth()

	// 1. Verificar disponibilidad básica de los endpoints
	endpoints := map[string]string{
		"signing":     config.MHPaths.AuthURL,
		"reception":   config.MHPaths.ReceptionURL,
		"contingency": config.MHPaths.ContingencyURL,
	}

	for name, url := range endpoints {
		if err := c.checkEndpoint(url); err != nil {
			logs.Error(fmt.Sprintf("Hacienda %s endpoint unavailable", name), map[string]interface{}{
				"error": err.Error(),
				"url":   url,
			})

			result.Down()
			result.AddInfo("error", fmt.Sprintf("Endpoint %s unavailable: %s", name, err.Error()))
			return result
		}
	}

	// 2. Verificar el procesamiento real mediante intento de autenticación
	if err := c.checkAuthProcessing(); err != nil {
		logs.Error("Hacienda signing processing check failed", map[string]interface{}{
			"error": err.Error(),
		})

		result.Down()
		result.AddInfo("error", fmt.Sprintf("Authentication processing failed: %s", err.Error()))
		return result
	}

	result.Up()
	return result
}

// Mantener los métodos auxiliares originales sin cambios
func (c *haciendaChecker) checkEndpoint(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	if resp.StatusCode >= 500 {
		return fmt.Errorf("service unavailable (status: %d)", resp.StatusCode)
	}

	return nil
}

func (c *haciendaChecker) checkAuthProcessing() error {
	// Credenciales de prueba
	dummyAuth := struct {
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	}{
		User: "test_user",
		Pwd:  "test_password",
	}

	jsonData, err := json.Marshal(dummyAuth)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		"POST",
		config.MHPaths.AuthURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	// Si el servicio está funcionando, debería responder con 401 o 400
	// ya que las credenciales son inválidas
	if resp.StatusCode != http.StatusUnauthorized &&
		resp.StatusCode != http.StatusBadRequest &&
		resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected service response: %d", resp.StatusCode)
	}

	return nil
}
