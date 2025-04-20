package ports

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
)

// FailedSequenceNumberRepositoryPort define la interfaz para el repositorio de números de secuencia fallidos
type FailedSequenceNumberRepositoryPort interface {
	// RegisterFailedSequence registra un número de secuencia fallido con detalles
	RegisterFailedSequence(
		ctx context.Context,
		branchID uint,
		dteType string,
		sequenceNumber uint,
		year uint,
		failureReason string,
		responseCode string,
		originalRequestData interface{},
		mhResponse string,
	) error

	// GetFailedSequences devuelve una lista de números de secuencia fallidos para una sucursal y tipo de DTE específicos
	GetFailedSequences(
		ctx context.Context,
		branchID uint,
		dteType string,
		limit int,
	) ([]db_models.FailedSequenceNumber, error)

	// GetFailedSequencesByYear devuelve una lista de números de secuencia fallidos para una sucursal, tipo de DTE y año específicos
	GetFailedSequencesByYear(
		ctx context.Context,
		branchID uint,
		dteType string,
		year uint,
		limit int,
	) ([]db_models.FailedSequenceNumber, error)
}
