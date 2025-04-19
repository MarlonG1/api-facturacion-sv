package dte_errors

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
)

type DTEError struct {
	ValidationErrors []error     // Errores de validación de value objects
	BusinessErrors   []*DTEError // Errores de reglas de negocio DTE
	ErrorType        string
	Message          string
	Code             string // Campo explícito para el código de error
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
		Code:             strings.ToUpper(errorType),
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

	return &DTEError{
		BusinessErrors: validErrors,
		Message:        strings.Join(messages, "; "),
		Code:           "MANY_ERRORS",
	}
}

// Error Implementación de la interfaz error para el error DTE
func (e *DTEError) Error() string {
	if e == nil {
		return "Unknown DTE error"
	}

	return e.Message
}

// GetValidationErrorsString Obtiene los errores de validación asociados al error DTE en caso de existir
func (e *DTEError) GetValidationErrorsString() []string {
	var messages []string
	for _, err := range e.ValidationErrors {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	if e.BusinessErrors != nil {
		for _, err := range e.BusinessErrors {
			if err != nil {
				messages = append(messages, err.Message)
			}
		}
	}

	return messages
}

// GetCode Obtiene el código de error
func (e *DTEError) GetCode() string {
	if e == nil {
		return "UNKNOWN_DTE_ERROR"
	}

	if e.Code != "" {
		return e.Code
	}

	return strings.ToUpper(e.ErrorType)
}

// GetMessage Obtiene el mensaje traducido del error
func (e *DTEError) GetMessage() string {
	if e == nil {
		return "Unknown DTE error"
	}

	if len(e.ValidationErrors) > 0 || len(e.BusinessErrors) > 0 {
		return i18n.Translate("service_errors.FailedToCreateDTE")
	}

	key := fmt.Sprintf("service_errors.%s", e.ErrorType)
	translated := i18n.Translate(key)

	// Si no hay traducción específica, usa el mensaje original
	if translated == key {
		return e.Message
	}

	return translated
}
