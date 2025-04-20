package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/validator/strategy"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type CCFRulesValidator struct {
	document   *ccf_models.CreditFiscalDocument
	strategies []interfaces.DTEValidationStrategy
}

// NewCCFRulesValidator Crea un validador de reglas para CCF
func NewCCFRulesValidator(doc *ccf_models.CreditFiscalDocument) *CCFRulesValidator {
	validator := &CCFRulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy.CCFItemStrategy{Document: doc},       // Validaciones de items
			&strategy.CCFTaxStrategy{Document: doc},        // Validaciones de impuestos específicos
			&strategy.CCFReceiverStrategy{Document: doc},   // Validaciones de receptor
			&strategy.CCFRelatedDocStrategy{Document: doc}, // Validaciones de documentos relacionados
		},
	}
	return validator
}

// Validate Ejecuta las validaciones de comprobante de crédito fiscal.
func (v *CCFRulesValidator) Validate() *dte_errors.DTEError {
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
