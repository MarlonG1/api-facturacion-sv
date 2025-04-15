package checkers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/dimiro1/health/db"
	"gorm.io/gorm"
)

type databaseChecker struct {
	db *gorm.DB
}

func NewDatabaseChecker(db *gorm.DB) health.ComponentChecker {
	return &databaseChecker{db: db}
}

func (c *databaseChecker) Name() string {
	return "database"
}

func (c *databaseChecker) Check() models.Health {
	// 1. Check the database connection
	sql, err := c.db.DB()
	if err != nil {
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Failed to get database connection",
		}
	}

	checker := db.NewMySQLChecker(sql)
	health := checker.Check()

	// 2. Check if the database is up
	if health.IsDown() {
		details := "Database connection is down"

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
		Details: "Database is healthy",
	}

}
