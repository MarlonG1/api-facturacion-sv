package ports

import (
	"context"
	authModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
)

// BatchTransmitterPort interfaz para transmisión de lotes a Hacienda
type BatchTransmitterPort interface {
	TransmitBatch(ctx context.Context, systemNIT string, dteType string, documents []string, token string, credentials authModels.HaciendaCredentials) (*models.BatchResponse, error)
	VerifyContingencyBatchStatus(ctx context.Context, batchID string, mhBatchID string, token string, docsMap map[string]dte.ContingencyDocument) error
	GetDTEVersion(dteType string) int
}
