package interfaces

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
)

// RetentionManager es una interfaz que define los métodos para la creación y validación de comprobantes de retención electrónicos
type RetentionManager interface {
	// Create crea una retención electrónica
	Create(ctx context.Context, data *retention_models.InputRetentionData, branchID uint, isAllPhysical bool) (*retention_models.RetentionModel, error)
	// Validate valida una retención electrónica
	Validate(retention *retention_models.RetentionModel) error
	// IsValid verifica si una retención electrónica  es válida
	IsValid(retention *retention_models.RetentionModel) bool
}
