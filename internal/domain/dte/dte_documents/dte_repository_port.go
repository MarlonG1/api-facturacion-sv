package dte_documents

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTERepositoryPort es una interfaz que define los métodos de un repositorio de DTE.
type DTERepositoryPort interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error
	// Update actualiza el estado de un DTE en la base de datos.
	Update(ctx context.Context, branchID uint, document dte.DTEDetails) error
	// GetByGenerationCode obtiene un DTE por su código de generación para consultas.
	GetByGenerationCode(ctx context.Context, branchID uint, id string) (*dte.DTEDocument, error)
	// GetDTEBalanceControl obtiene el control de saldo de un DTE por su ID.
	GetDTEBalanceControl(ctx context.Context, branchID uint, id string) (*dte.BalanceControl, error)
	// GenerateBalanceTransaction genera una transacción de saldo para un DTE.
	GenerateBalanceTransaction(ctx context.Context, branchID uint, originalDTE string, transaction *dte.BalanceTransaction) error
	// VerifyStatus verifica el estado de un DTE en la base de datos.
	VerifyStatus(ctx context.Context, branchID uint, id string) (string, error)
	// GetTotalCount obtiene el número total de DTEs en la base de datos.
	GetTotalCount(ctx context.Context, filters *dte.DTEFilters) (int64, error)
	// GetSummaryStats obtiene las estadísticas resumidas de los DTEs en la base de datos.
	GetSummaryStats(ctx context.Context, filters *dte.DTEFilters) (*dte.ListSummary, error)
	// GetPagedDocuments obtiene una lista paginada de DTEs en la base de datos.
	GetPagedDocuments(ctx context.Context, filters *dte.DTEFilters) ([]dte.DTEModelResponse, error)
}
