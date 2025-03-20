package interfaces

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
)

type CCFManager interface {
	// Operaciones principales
	Create(data *ccf_models.CCFData) (*ccf_models.CreditFiscalDocument, error) // Create crea una invoice electrónica a partir de los datos de una invoice
	Validate(invoice *ccf_models.CreditFiscalDocument) error                   // Validate valida una invoice electrónica

	// Generación de identificadores
	GenerateControlNumber(invoice *ccf_models.CreditFiscalDocument) error      // GenerateControlNumber genera el número de control de la invoice electrónica
	GenerateCodeAndIdentifiers(invoice *ccf_models.CreditFiscalDocument) error // GenerateCodeAndIdentifiers genera el código de control y los identificadores de la invoice electrónica

	// Estado y validación
	IsValid(invoice *ccf_models.CreditFiscalDocument) bool // IsValid verifica si una invoice electrónica es válida
}
