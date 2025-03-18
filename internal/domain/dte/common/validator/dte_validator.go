package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

// ValidateDTEDocument valida un documento DTE y retorna un error si no cumple con las reglas de negocio
// Se manejan 3 tipos de errores:
// 1. Errores de validación de DTE (errores de validación de reglas de negocio DTE)
// 2. Errores de validación de value objects
// 3. Errores de validación de value objects y DTE, retorna un CompositeError con los errores de validación
func ValidateDTEDocument[T interfaces.DTEValidator](doc T) error {
	//Tipo 1 - Errores de validación de DTE
	validationErrors := ValidateModel(doc)

	//Tipo 2 - Errores de validación de value objects
	if dteErr := doc.ValidateDTERules(); dteErr != nil {
		logs.Error("Failed to validate DTE rules", map[string]interface{}{"error": dteErr.Error()})
		dteErr.ValidationErrors = append(dteErr.ValidationErrors, validationErrors...)
		return dteErr
	}

	//Tipo 3 - Errores de validación de value objects y DTE
	if len(validationErrors) > 0 {
		logs.Error("Failed to validate DTE document", map[string]interface{}{"error": dte_errors.NewCompositeError(validationErrors...).Error()})
		return dte_errors.NewCompositeError(validationErrors...)
	}

	return nil
}
