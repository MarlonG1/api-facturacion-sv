package interfaces

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
)

// CCFManager es la interfaz que define las operaciones de un manager de Comprobante de Crédito Fiscal
type CCFManager interface {
	// Create crea una invoice electrónica a partir de los datos de una invoice
	Create(ctx context.Context, data *ccf_models.CCFData, branchID uint) (*ccf_models.CreditFiscalDocument, error)
	// Validate valida una invoice electrónica
	Validate(invoice *ccf_models.CreditFiscalDocument) error
	// IsValid verifica si una invoice electrónica es válida
	IsValid(invoice *ccf_models.CreditFiscalDocument) bool
}
