package checkers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/dimiro1/health/url"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
)

type signerChecker struct {
	client *http.Client
}

func NewSignerChecker() ports.ComponentChecker {
	return &signerChecker{
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (s *signerChecker) Name() string {
	return "dte_signer"
}

func (s *signerChecker) Check() models.Health {
	checker := url.NewCheckerWithTimeout(config.Signer.Health, s.client.Timeout)
	health := checker.Check()

	if health.IsDown() {
		details := "Signer service is down"
		if health.GetInfo("error") != nil {
			details = fmt.Sprintf("%s: %v", details, health.GetInfo("error"))
		}

		return models.Health{
			Status:  constants.StatusDown,
			Details: details,
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: "Signer service is healthy",
	}
}
