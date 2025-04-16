package location

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Address struct {
	Value string `json:"value"`
}

func NewAddress(value string) (*Address, error) {
	addr := &Address{Value: value}
	if addr.IsValid() {
		return addr, nil
	}
	return &Address{}, dte_errors.NewValidationError("InvalidLength", "address", "1-200", value)
}

func NewValidatedAddress(value string) *Address {
	return &Address{Value: value}
}

// IsValid verifica si el valor de la dirección es válido (1 a 200 caracteres)
func (a *Address) IsValid() bool {
	return len(a.Value) >= 1 && len(a.Value) <= 200
}

func (a *Address) Equals(other interfaces.ValueObject[string]) bool {
	return a.GetValue() == other.GetValue()
}

func (a *Address) GetValue() string {
	return a.Value
}

func (a *Address) ToString() string {
	return a.Value
}
