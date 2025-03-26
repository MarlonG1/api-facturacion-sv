package interfaces

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTEManager es una interfaz que define los métodos de un administrador de DTE.
type DTEManager interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(context.Context, interface{}, string, string, *string) error
	// UpdateDTE actualiza el estado de un DTE en la base de datos.
	UpdateDTE(ctx context.Context, id, status string, receptionStamp *string) error
	// GetByGenerationCode obtiene un DTE por su código de generación para consultas.
	GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEResponse, error)
}
