package document

import (
	"reflect"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type DTEType struct {
	Value string `json:"value"`
}

// NewDTEType crea un nuevo tipo de documento electrónico válido para emisión
func NewDTEType(value string) (*DTEType, error) {
	tipoDte := &DTEType{Value: value}
	if tipoDte.IsValid() {
		return tipoDte, nil
	}
	return &DTEType{}, dte_errors.NewValidationError("InvalidFormat", "tipoDte", "un tipo de documento válido", value)
}

func NewValidatedDTEType(value string) *DTEType {
	return &DTEType{Value: value}
}

// NewDTETypeForReceiver crea un nuevo tipo de documento electrónico válido para recepción
func NewDTETypeForReceiver(value string) (*DTEType, error) {
	tipoDte := &DTEType{Value: value}
	if tipoDte.IsForReception() {
		return tipoDte, nil
	}
	return &DTEType{}, dte_errors.NewValidationError("InvalidFormat", "tipoDte", "un tipo de documento válido para recepción", value)
}

// IsValid válido si el valor es un string y es un tipo de documento electrónico válido
func (t *DTEType) IsValid() bool {
	for _, v := range constants.ValidDTETypes {
		if t.Value == v {
			return true
		}
	}
	return false
}

func (t *DTEType) IsForReception() bool {
	for _, v := range constants.ValidReceiverDTETypes {
		if t.Value == v {
			return true
		}
	}
	return false
}

func (t *DTEType) Equals(other interfaces.ValueObject[string]) bool {
	return t.GetValue() == other.GetValue()
}

func (t *DTEType) ToString() string {
	return reflect.ValueOf(t.Value).String()
}

func (t *DTEType) GetValue() string {
	return t.Value
}
