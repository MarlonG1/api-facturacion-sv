package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type DeliveryName struct {
	Value string `json:"value"`
}

func NewDeliveryName(value string) (*DeliveryName, error) {
	name := &DeliveryName{Value: value}
	if name.IsValid() {
		return name, nil
	}
	return &DeliveryName{}, dte_errors.NewValidationError("InvalidLength", "DeliveryName", "1 a 100", value)
}

func NewValidatedDeliveryName(value string) *DeliveryName {
	return &DeliveryName{Value: value}
}

func (n *DeliveryName) IsValid() bool {
	return len(n.Value) >= 1 && len(n.Value) <= 100
}

func (n *DeliveryName) Equals(other interfaces.ValueObject[string]) bool {
	return n.GetValue() == other.GetValue()
}

func (n *DeliveryName) GetValue() string {
	return n.Value
}

func (n *DeliveryName) ToString() string {
	return n.Value
}
