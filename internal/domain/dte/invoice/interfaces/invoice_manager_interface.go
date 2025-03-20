package interfaces

import (
	invoice_models2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
)

type InvoiceManager interface {
	// Operaciones principales
	Create(data *invoice_models2.InvoiceData) (*invoice_models2.ElectronicInvoice, error) // Create crea una invoice electrónica a partir de los datos de una invoice
	Validate(invoice *invoice_models2.ElectronicInvoice) error                            // Validate valida una invoice electrónica

	// Generación de identificadores
	GenerateControlNumber(invoice *invoice_models2.ElectronicInvoice) error      // GenerateControlNumber genera el número de control de la invoice electrónica
	GenerateCodeAndIdentifiers(invoice *invoice_models2.ElectronicInvoice) error // GenerateCodeAndIdentifiers genera el código de control y los identificadores de la invoice electrónica

	// Estado y validación
	IsValid(invoice *invoice_models2.ElectronicInvoice) bool // IsValid verifica si una invoice electrónica es válida
}
