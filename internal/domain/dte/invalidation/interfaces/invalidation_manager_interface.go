package interfaces

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
)

type InvalidationManager interface {
	// Validación y creación
	Create(data *models.InvalidationDocument) (*models.InvalidationDocument, error)
	Validate(document *models.InvalidationDocument) error
	InvalidateDocument(ctx context.Context, generationCode string) error

	// Estado y validación
	IsValid(document *models.InvalidationDocument) bool
	ValidateDTEExists(document *models.InvalidationDocument) error
}
