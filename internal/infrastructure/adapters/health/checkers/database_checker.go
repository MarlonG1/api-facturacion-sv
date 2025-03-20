package checkers

import (
	"gorm.io/gorm"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type databaseChecker struct {
	db *gorm.DB
}

func NewDatabaseChecker(db *gorm.DB) ports.ComponentChecker {
	return &databaseChecker{db: db}
}

func (c *databaseChecker) Name() string {
	return "database"
}

func (c *databaseChecker) Check() models.Health {
	db, err := c.db.DB()
	if err != nil {
		logs.Error("Database connection error", map[string]interface{}{
			"error": err.Error(),
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Database service unavailable",
		}
	}

	if err := db.Ping(); err != nil {
		logs.Error("Database ping failed", map[string]interface{}{
			"error": err.Error(),
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Database service ping failed",
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: "Database service available",
	}
}
