package health

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/health/checkers"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"gorm.io/gorm"
)

type healthService struct {
	checkers []ports.ComponentChecker
}

type HealthServiceConfig struct {
	DB *gorm.DB
}

func NewHealthService(cfg *HealthServiceConfig) ports.HealthManager {
	service := &healthService{
		checkers: []ports.ComponentChecker{
			checkers.NewDatabaseChecker(cfg.DB),
			checkers.NewRedisChecker(),
			checkers.NewHaciendaChecker(),
			checkers.NewFileSystemChecker(),
			checkers.NewSignerChecker(),
		},
	}
	return service
}

func (s *healthService) CheckHealth() (*models.HealthStatus, error) {
	components := make(map[string]models.Health)
	status := constants.StatusUp

	for _, checker := range s.checkers {
		health := checker.Check()
		components[checker.Name()] = health

		if health.Status == constants.StatusDown {
			status = constants.StatusDown
		}
	}

	return &models.HealthStatus{
		Status:     status,
		Components: components,
		Timestamp:  utils.TimeNow().Format("02-01-2006 15:04:05"),
	}, nil
}
