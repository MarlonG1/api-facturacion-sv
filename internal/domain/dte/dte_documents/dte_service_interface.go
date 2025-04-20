package dte_documents

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTEManager es una interfaz que define los métodos de un administrador de DTE.
type DTEManager interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(context.Context, interface{}, string, string, *string) error
	// UpdateDTE actualiza el estado de un DTE en la base de datos.
	UpdateDTE(ctx context.Context, branchID uint, document dte.DTEDetails) error
	// VerifyStatus verifica el estado de un DTE en la base de datos.
	VerifyStatus(ctx context.Context, branchID uint, id string) (string, error)
	// GetByGenerationCode obtiene un DTE por su código de generación para procesos internos.
	GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error)
	// GenerateBalanceTransaction genera una transacción de balance para un DTE.
	GenerateBalanceTransaction(ctx context.Context, branchID uint, transactionType, id, originalDTE string, document interface{}) error
	// ValidateForCreditNote valida un DTE para la creación de una Nota de Crédito.
	ValidateForCreditNote(ctx context.Context, branchID uint, originalDTE string, document interface{}) error
	// GetByGenerationCodeConsult obtiene un DTE por su código de generación para consultas.
	GetByGenerationCodeConsult(ctx context.Context, branchID uint, generationCode string) (*dte.DTEResponse, error)
	// GetAllDTEs obtiene todos los DTEs en la base de datos con filtros y paginación.
	GetAllDTEs(ctx context.Context, filters *dte.DTEFilters) (*dte.DTEListResponse, error)
}
