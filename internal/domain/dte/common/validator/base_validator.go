package validator

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"reflect"
)

// ValidateModel valida un modelo y retorna los errores de validación encontrados
// Se valida que los campos de tipo ValueObject cumplan con las reglas de negocio definidas
func ValidateModel[T any](model T) []error {
	var validationErrors []error
	v := reflect.ValueOf(model)

	// Si el modelo es un puntero, se obtiene el valor apuntado
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return validationErrors
		}
		v = v.Elem()
	}

	// Si el valor no es una estructura, no se realiza la validación
	if v.Kind() != reflect.Struct {
		return validationErrors
	}

	// Se recorren los campos de la estructura para validar los campos de tipo ValueObject y ValueObject[]
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Skip si el campo no es exportado
		if !fieldType.IsExported() {
			continue
		}

		if field.Kind() == reflect.Slice {
			sliceErrors := validateSlice(field)
			validationErrors = append(validationErrors, sliceErrors...)
			continue
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		// Verificar que el campo sea válido antes de usarlo
		if !field.IsValid() {
			continue
		}

		// Si el campo es una estructura, se valida recursivamente
		if field.Kind() == reflect.Struct {
			// Verificar que podemos llamar a Interface()
			if field.CanInterface() {
				structErrors := ValidateModel(field.Interface())
				validationErrors = append(validationErrors, structErrors...)
			}
			continue
		}

		// Verificar que podemos llamar a Interface() antes de intentar la conversión
		if field.CanInterface() {
			if validator, ok := field.Interface().(interfaces.ValueObject[any]); ok {
				if !validator.IsValid() {
					validationErrors = append(validationErrors,
						dte_errors.NewValidationError("InvalidField", fieldType.Name))
				}
			}
		}
	}

	return validationErrors
}

// validateSlice valida los elementos de un slice y retorna los errores de validación encontrados
// Se valida que los campos de tipo ValueObject cumplan con las reglas de negocio definidas
func validateSlice(field reflect.Value) []error {
	var sliceErrors []error

	for i := 0; i < field.Len(); i++ {
		element := field.Index(i)

		if element.Kind() == reflect.Ptr {
			if element.IsNil() {
				continue
			}
			element = element.Elem()
		}

		// Verificar que el elemento sea válido
		if !element.IsValid() {
			continue
		}

		if element.Kind() == reflect.Struct {
			// Verificar que podemos llamar a Interface()
			if element.CanInterface() {
				elementErrors := ValidateModel(element.Interface())
				sliceErrors = append(sliceErrors, elementErrors...)
			}
		}
	}

	return sliceErrors
}
