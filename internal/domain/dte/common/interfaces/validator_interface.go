package interfaces

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"

// Validator Interfaz que define los métodos que deben ser implementados por los objetos que validan un campo
type Validator interface {
	IsValid() bool // IsValid Valida que el valor del campo cumpla con las reglas de negocio definidas
}

// DTEValidator Interfaz que define los métodos que deben ser implementados por los objetos que validan un DTE
type DTEValidator interface {
	ValidateDTERules() *dte_errors.DTEError // ValidateDTERules Valida las reglas de negocio de un DTE
}

// ValueObject Interfaz que define los métodos que deben ser implementados por los objetos de valor
type ValueObject[T any] interface {
	Validator                         // Validator Interfaz que define los métodos que deben ser implementados por los objetos que validan un DTE
	ToString() string                 // ToString Convierte el valor del campo a un string
	Equals(value ValueObject[T]) bool // Equals Compara el valor del campo con el valor de otro objeto de valor, a través de generics se define el tipo de dato que se espera
	GetValue() T                      // GetValue Obtiene el valor del campo, a través de generics se define el tipo de dato que se espera
}

// DTEValidationStrategy Interfaz que define los métodos que deben ser implementados por las estrategias de validación de DTE
type DTEValidationStrategy interface {
	Validate() *dte_errors.DTEError // Validate Valida un DTE
}
