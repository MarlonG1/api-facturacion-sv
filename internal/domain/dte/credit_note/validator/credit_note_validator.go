package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/validator/strategy"
)

type CreditNoteRulesValidator struct {
	document   *credit_note_models.CreditNoteModel
	strategies []interfaces.DTEValidationStrategy
}

func NewCreditNoteRulesValidator(doc *credit_note_models.CreditNoteModel) *CreditNoteRulesValidator {
	validator := &CreditNoteRulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy.CreditNoteItemStrategy{Document: doc},       // Validaciones de ítems
			&strategy.CreditNoteTaxStrategy{Document: doc},        // Validaciones de impuestos
			&strategy.CreditNoteRelatedDocStrategy{Document: doc}, // Validaciones de documentos relacionados
		},
	}
	return validator
}

// Validate Ejecuta las validaciones de la nota de crédito electrónica.
func (v *CreditNoteRulesValidator) Validate() *dte_errors.DTEError {
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
