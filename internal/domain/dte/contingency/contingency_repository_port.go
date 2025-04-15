package contingency

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"time"
)

// ContingencyRepositoryPort interfaz para el repositorio de contingencia (ya existe)
type ContingencyRepositoryPort interface {
	// Create almacena un documento de contingencia en la base de datos
	Create(ctx context.Context, doc *dte.ContingencyDocument) error
	// GetPending obtiene los documentos en estado PENDING para procesar
	GetPending(ctx context.Context, limit int) ([]dte.ContingencyDocument, error)
	// UpdateBatch actualiza el estado de los documentos de un lote
	UpdateBatch(ctx context.Context, ids []string, observations []string, stamps map[string]string, batchID string, mhBatchID string, status string) error
	// GetFirstContingencyTimestamp obtiene la fecha de la primera contingencia de un sistema
	GetFirstContingencyTimestamp(ctx context.Context, branchID uint) (*time.Time, error)
}
