package interfaces

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
)

// InvoiceManager es una interfaz que define los métodos para la creación y validación de facturas electrónicas
type InvoiceManager interface {
	// Create crea una invoice electrónica a partir de los datos de una invoice
	Create(ctx context.Context, data *invoice_models.InvoiceData, branchID uint) (*invoice_models.ElectronicInvoice, error)
	// Validate valida una invoice electrónica
	Validate(invoice *invoice_models.ElectronicInvoice) error
	// IsValid verifica si una invoice electrónica es válida
	IsValid(invoice *invoice_models.ElectronicInvoice) bool
}
