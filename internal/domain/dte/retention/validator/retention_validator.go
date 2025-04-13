package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/validator/strategy"
)

type RetentionRulesValidator struct {
	document   *retention_models.RetentionModel
	strategies []interfaces.DTEValidationStrategy
}

// NewRetentionRulesValidator Crea un validador de reglas para facturas electrónicas
func NewRetentionRulesValidator(doc *retention_models.RetentionModel) *RetentionRulesValidator {
	validator := &RetentionRulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy.RetentionItemStrategy{Document: doc},  // 1. Validaciones de items
			&strategy.RetentionTotalStrategy{Document: doc}, // 2. Validaciones de totales
		},
	}
	return validator
}

// Validate Ejecuta las validaciones de la invoice electrónica.
func (v *RetentionRulesValidator) Validate() *dte_errors.DTEError {
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
