package checkers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/dimiro1/health/url"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
)

type signerChecker struct {
	client *http.Client
}

func NewSignerChecker() health.ComponentChecker {
	return &signerChecker{
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *signerChecker) Name() string {
	return "dte_signer"
}

func (c *signerChecker) Check() models.Health {
	checker := url.NewCheckerWithTimeout(config.Signer.Health, c.client.Timeout)
	health := checker.Check()

	if health.IsDown() {
		details := utils.TranslateHealthDown(c.Name())
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
		Details: utils.TranslateHealthUp(c.Name()),
	}
}
