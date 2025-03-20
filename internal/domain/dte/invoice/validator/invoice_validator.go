package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	strategy2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/validator/strategy"
)

type InvoiceRulesValidator struct {
	document   *invoice_models.ElectronicInvoice
	strategies []interfaces.DTEValidationStrategy
}

// NewInvoiceRulesValidator Crea un validador de reglas para facturas electrónicas
func NewInvoiceRulesValidator(doc *invoice_models.ElectronicInvoice) *InvoiceRulesValidator {
	validator := &InvoiceRulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy2.InvoiceItemsStrategy{Document: doc},  // 1. Validaciones de items
			&strategy2.InvoiceTaxStrategy{Document: doc},    // 2. Validaciones de impuestos específicos
			&strategy2.InvoiceTotalsStrategy{Document: doc}, // 3. Cálculos específicos
		},
	}
	return validator
}

// Validate Ejecuta las validaciones de la invoice electrónica.
func (v *InvoiceRulesValidator) Validate() *dte_errors.DTEError {
	var validationErrors []*dte_errors.DTEError

	for _, strategyValidator := range v.strategies {
		if err := strategyValidator.Validate(); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return dte_errors.NewDTEErrorComposite(validationErrors)
	}

	return nil
}
