package dte_errors

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
)

type ValidationError struct {
	ErrorType string
	Message   string
}

// NewValidationError Crea un nuevo error de validación con el tipo de error y los parámetros enviados
func NewValidationError(errorType string, params ...interface{}) *ValidationError {
	message := constants.GetErrorMessage(errorType, params...)
	return &ValidationError{ErrorType: errorType, Message: message}
}

func NewFormattedValidationError(err error) *ValidationError {
	return &ValidationError{ErrorType: "", Message: err.Error()}
}

// Error Implementación de la interfaz error para el error de validación
func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s", v.Message)
}

// GetType Retorna el tipo de error de validación
func (v *ValidationError) GetType() string {
	if v == nil {
		return "UnknownError"
	}

	return v.ErrorType
}
