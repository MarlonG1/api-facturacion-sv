package invalidation

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// InvalidationManager es la interfaz que define los métodos que se pueden realizar sobre la invalidación de documentos
type InvalidationManager interface {
	// Validate valida el documento de invalidación
	Validate(ctx context.Context, branchID uint, document *models.InvalidationDocument) error
	// ValidateStatus valida el estado del documento a invalidar y del documento de reemplazo
	ValidateStatus(ctx context.Context, branchID uint, req structs.InvalidationRequest) error
	// InvalidateDocument invalida un documento
	InvalidateDocument(ctx context.Context, branchID uint, originalCode string) error
}
