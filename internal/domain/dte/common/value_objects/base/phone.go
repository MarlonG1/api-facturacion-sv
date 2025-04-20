package base

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Phone struct {
	Value string `json:"value"`
}

func NewPhone(value string) (*Phone, error) {
	phone := &Phone{Value: value}
	if phone.IsValid() {
		return phone, nil
	}
	return &Phone{}, dte_errors.NewValidationError("InvalidPhone", value)
}

func NewValidatedPhone(value string) *Phone {
	return &Phone{Value: value}
}

// IsValid válida si el número de teléfono tiene entre 8 y 30 dígitos
func (p *Phone) IsValid() bool {
	return len(p.Value) >= 8 && len(p.Value) <= 30
}

func (p *Phone) Equals(other interfaces.ValueObject[string]) bool {
	return p.GetValue() == other.GetValue()
}

func (p *Phone) GetValue() string {
	return p.Value
}

func (p *Phone) ToString() string {
	return p.Value
}
