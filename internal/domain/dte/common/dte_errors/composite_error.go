package dte_errors

import "strings"

// CompositeError representa un error compuesto que contiene varios errores
type CompositeError struct {
	Errors []error
}

func NewCompositeError(errors ...error) *CompositeError {
	return &CompositeError{
		Errors: errors,
	}
}

// Error Implementación de la interfaz error para el error compuesto para lanzar error principal a nivel de DTEDocument y errores de validación de value objects
func (e *CompositeError) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}
