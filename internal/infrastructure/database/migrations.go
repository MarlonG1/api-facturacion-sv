package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

// modelsToMigrate contiene todos los modelos que se deben migrar
var modelsToMigrate = []schema.Tabler{
	&db_models.User{},
	&db_models.BranchOffice{},
	&db_models.Address{},
	&db_models.DTEDetails{},
	&db_models.DTEDocument{},
	&db_models.ContingencyDocument{},
	&db_models.ControlNumberSequence{},
	&db_models.SystemAdmin{},
	&db_models.DomainEvent{},
	&db_models.UserNotification{},
	&db_models.AdminAlert{},
	&db_models.AdminAlertRecipients{},
	&db_models.NotifiableUser{},
}

// RunMigrations ejecuta todas las migraciones de la base de datos
func RunMigrations(db *gorm.DB) error {
	logs.Info("Starting database migrations")

	for i, model := range modelsToMigrate {
		tn := model.TableName()
		logs.Info(fmt.Sprintf("Starting model migration #%d: %s", i, tn))

		if err := db.AutoMigrate(model); err != nil {
			logs.Error("Failed to migrate model", map[string]interface{}{
				"index": i,
				"model": tn,
				"error": err.Error(),
			})
			return err
		}

		logs.Info(fmt.Sprintf("Successfully migrated model %s", tn))
	}

	logs.Info("All migrations completed successfully")
	return nil
}
