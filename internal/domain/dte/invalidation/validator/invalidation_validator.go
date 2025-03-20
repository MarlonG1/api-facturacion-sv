package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/validator/strategy"
)

type InvalidationRulesValidator struct {
	document   *models.InvalidationDocument
	strategies []interfaces.DTEValidationStrategy
}

func NewInvalidationRulesValidator(doc *models.InvalidationDocument) *InvalidationRulesValidator {
	validator := &InvalidationRulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy.InvalidationBasicStrategy{Document: doc},    // Validaciones básicas
			&strategy.InvalidationDocumentStrategy{Document: doc}, // Validaciones del documento a invalidar
			&strategy.InvalidationReasonStrategy{Document: doc},   // Validaciones específicas por tipo de anulación
			&strategy.InvalidationDateStrategy{Document: doc},     // Validaciones de fechas y plazos
		},
	}
	return validator
}

func (v *InvalidationRulesValidator) Validate() *dte_errors.DTEError {
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
