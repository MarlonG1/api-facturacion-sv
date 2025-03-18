package item

import (
	"fmt"
	"math"
	"strconv"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Quantity struct {
	Value float64 `json:"value"`
}

func NewQuantity(value float64) (*Quantity, error) {
	quantity := &Quantity{Value: value}
	if quantity.IsValid() {
		return quantity, nil
	}
	return &Quantity{}, dte_errors.NewValidationError("InvalidQuantity", fmt.Sprintf("%f", value))
}

func NewValidatedQuantity(value float64) *Quantity {
	return &Quantity{Value: value}
}

// IsValid válida que el valor de Quantity sea mayor o igual a 0 y menor o igual a 99999999999.99
func (a *Quantity) IsValid() bool {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", a.Value), 64)
	return value > 0 && value <= 99999999999.99
}

func (a *Quantity) Equals(other interfaces.ValueObject[float64]) bool {
	return math.Abs(a.GetValue()-other.GetValue()) < 0.000001
}

func (a *Quantity) GetValue() float64 {
	return a.Value
}

func (a *Quantity) ToString() string {
	return fmt.Sprintf("%.2f", a.Value)
}
