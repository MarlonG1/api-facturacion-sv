package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator/strategy"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type DTERulesValidator struct {
	document   interfaces.DTEDocument
	strategies []interfaces.DTEValidationStrategy
}

// NewDTERulesValidator Crea un validador de reglas de DTE
func NewDTERulesValidator(doc interfaces.DTEDocument) *DTERulesValidator {
	validator := &DTERulesValidator{
		document: doc,
		strategies: []interfaces.DTEValidationStrategy{
			&strategy.BasicRulesStrategy{Document: doc},         // 1. Validaciones básicas
			&strategy.TemporalValidationStrategy{Document: doc}, // 2. Validaciones temporales
			&strategy.ContingencyStrategy{Document: doc},        // 3. Validaciones de contingencia
			&strategy.ModelTypeStrategy{Document: doc},          // 4. Validaciones de tipo de modelo
			&strategy.TaxCalculationStrategy{Document: doc},     // 5. Validaciones de cálculo de impuestos
			&strategy.PaymentTotalStrategy{Document: doc},       // 6. Validaciones de total de pagos
			&strategy.ItemValidationStrategy{Document: doc},     // 7. Validaciones de items
			&strategy.ExtensionStrategy{Document: doc},          // 8. Validaciones de extensión
			&strategy.DocumentTypeStrategy{Document: doc},       // 9. Validaciones de tipo de documento
			&strategy.RelatedDocsStrategy{Document: doc},        // 10. Validaciones de documentos relacionados
			&strategy.ThirdPartyStrategy{Document: doc},         // 11. Validaciones de venta a terceros
			&strategy.OtherDocumentsStrategy{Document: doc},     // 12. Validaciones de otros documentos
		},
	}
	return validator
}

// Validate Válida las reglas de un documento DTE según las estrategias definidas
func (v *DTERulesValidator) Validate() *dte_errors.DTEError {
	var validationErrors []*dte_errors.DTEError

	for i, strategyValidator := range v.strategies {
		logs.Info("Starting validation for strategy", map[string]interface{}{"strategy": i})
		if err := strategyValidator.Validate(); err != nil {
			logs.Error("Failed to validate strategy", map[string]interface{}{"strategy": i, "error": err.Error()})
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		// Crear un DTEError compuesto
		return dte_errors.NewDTEErrorComposite(validationErrors)
	}

	return nil
}
