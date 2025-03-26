package ports

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTERepositoryPort es una interfaz que define los métodos de un repositorio de DTE.
type DTERepositoryPort interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error
	// Update actualiza el estado de un DTE en la base de datos.
	Update(ctx context.Context, id, status string, receptionStamp *string) error
	// GetByGenerationCode obtiene un DTE por su código de generación para consultas.
	GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error)
}
