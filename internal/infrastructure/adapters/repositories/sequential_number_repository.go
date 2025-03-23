package repositories

import (
	"context"
	"errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"gorm.io/gorm"
)

type ControlNumberRepository struct {
	db *gorm.DB
}

// NewControlNumberRepository crea una instancia de ControlNumberRepository. Recibe una instancia de gorm.DB.
func NewControlNumberRepository(db *gorm.DB) ports.SequentialNumberRepositoryPort {
	return &ControlNumberRepository{db: db}
}

// GetNext obtiene el siguiente número de control para un tipo de DTE, NIT de sistema y código de establecimiento.
func (r *ControlNumberRepository) GetNext(ctx context.Context, dteType string, branchID uint) (int, error) {
	currentYear := utils.TimeNow().Year()
	var sequence db_models.ControlNumberSequence

	// 1. Crear transacción para obtener el siguiente número de control de la secuencia
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// FOR UPDATE bloquea la fila para otras transacciones
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("branch_id = ? AND dte_type = ? AND year = ?", branchID, dteType, currentYear).
			First(&sequence)

		// 1.1 Si no existe la secuencia, crear una nueva
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			sequence = db_models.ControlNumberSequence{
				BranchID:   branchID,
				DTEType:    dteType,
				LastNumber: 0,
				Year:       currentYear,
			}
		}

		// 2. Incrementar secuencia
		sequence.LastNumber++

		// 3. Guardar o actualizar secuencia
		if result.Error == nil {
			// 3.1 Actualizar secuencia si existe
			if err := tx.Save(&sequence).Error; err != nil {
				return err
			}
		} else {
			// 3.2 Crear secuencia si no existe
			if err := tx.Create(&sequence).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logs.Error("Failed to get next control number", map[string]interface{}{
			"dteType":  dteType,
			"branchID": branchID,
			"error":    err.Error(),
		})
		return 0, err
	}

	return sequence.LastNumber, nil
}
