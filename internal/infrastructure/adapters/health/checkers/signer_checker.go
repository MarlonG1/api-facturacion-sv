package checkers

import (
	"context"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type signerChecker struct {
	client *http.Client
}

func NewSignerChecker() ports.ComponentChecker {
	return &signerChecker{
		client: &http.Client{Timeout: 1 * time.Second},
	}
}

func (s *signerChecker) Name() string {
	return "dte_signer"
}

func (s *signerChecker) Check() models.Health {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", env.Signer.Health, nil)
	if err != nil {
		logs.Error("Error creating signer service request", map[string]interface{}{
			"error": err.Error(),
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Service unavailable",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		logs.Error("Error calling signer service", map[string]interface{}{
			"error": err.Error(),
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Service unavailable",
		}
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logs.Error("Signer service unavailable", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Service unavailable",
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: "Service available",
	}
}
