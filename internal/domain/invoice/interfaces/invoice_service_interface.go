package interfaces

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/invoice/invoice_models"
)

type InvoiceManager interface {
	// Operaciones principales
	Create(data *invoice_models.InvoiceData) (*invoice_models.ElectronicInvoice, error) // Create crea una invoice electrónica a partir de los datos de una invoice
	Validate(invoice *invoice_models.ElectronicInvoice) error                           // Validate valida una invoice electrónica

	// Generación de identificadores
	GenerateControlNumber(invoice *invoice_models.ElectronicInvoice) error      // GenerateControlNumber genera el número de control de la invoice electrónica
	GenerateCodeAndIdentifiers(invoice *invoice_models.ElectronicInvoice) error // GenerateCodeAndIdentifiers genera el código de control y los identificadores de la invoice electrónica

	// Estado y validación
	IsValid(invoice *invoice_models.ElectronicInvoice) bool // IsValid verifica si una invoice electrónica es válida
}
