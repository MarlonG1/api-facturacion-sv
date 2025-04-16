package dte_errors

import (
	"fmt"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
)

type DTEError struct {
	ValidationErrors []error     // Errores de validación de value objects
	BusinessErrors   []*DTEError // Errores de reglas de negocio DTE
	ErrorType        string
	Message          string
}

// getDTEErrorMessage Obtiene el mensaje de error DTE con los parámetros enviados
func getDTEErrorMessage(errorType string, params ...interface{}) string {
	return constants.GetErrorMessage(errorType, params...)
}

// NewDTEErrorSimple Crea un nuevo error DTE con el tipo de error y los parámetros enviados sin errores de validación
func NewDTEErrorSimple(errorType string, params ...interface{}) *DTEError {
	return &DTEError{
		ValidationErrors: nil,
		ErrorType:        errorType,
		Message:          getDTEErrorMessage(errorType, params...),
	}
}

// NewDTEErrorComposite Crea un nuevo error DTE con los errores de negocio enviados y los agrupa en un solo error DTE
func NewDTEErrorComposite(businessErrors []*DTEError) *DTEError {
	var messages []string
	var validErrors []*DTEError

	for _, err := range businessErrors {
		if err != nil {
			validErrors = append(validErrors, err)
			messages = append(messages, err.Error())
		}
	}

	var errorType string
	if len(validErrors) > 0 {
		errorType = validErrors[0].ErrorType
	}

	return &DTEError{
		BusinessErrors: validErrors,
		ErrorType:      errorType,
		Message:        strings.Join(messages, "; "),
	}
}

// Error Implementación de la interfaz error para el error DTE para lanzar error principal a nivel de DTEDocument y errores de validación
func (e *DTEError) Error() string {
	if e == nil {
		return "unknown DTEError"
	}

	var messages []string
	messages = append(messages, fmt.Sprintf("%s", e.Message))

	if len(e.ValidationErrors) > 0 {
		messages = append(messages, "Errors:")
		for _, err := range e.ValidationErrors {
			messages = append(messages, fmt.Sprintf("- %s", err.Error()))
		}
	}

	return strings.Join(messages, " | -> | ")
}

// GetValidationErrorsString Obtiene los errores de validación asociados al error DTE en caso de existir
func (e *DTEError) GetValidationErrorsString() []string {
	var messages []string
	for _, err := range e.ValidationErrors {
		messages = append(messages, err.Error())
	}

	return messages
}
