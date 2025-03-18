package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

// TaxAmount es una estructura que representa un monto de un impuesto de un DTE, contiene TotalAmount
type TaxAmount struct {
	TotalAmount financial.Amount `json:"totalAmount,omitempty"`
}

func (t *TaxAmount) GetTotalAmount() financial.Amount {
	return t.TotalAmount
}

// Tax es una estructura que representa un impuesto de un DTE, contiene Code, Description y Value
type Tax struct {
	Code        financial.TaxType `json:"code"`
	Description string            `json:"description"`
	Value       *TaxAmount        `json:"value,omitempty"`
}

func (t *Tax) GetTotalAmount() float64 {
	return t.Value.TotalAmount.GetValue()
}

func (t *Tax) GetCode() string {
	return t.Code.GetValue()
}
func (t *Tax) GetDescription() string {
	return t.Description
}
func (t *Tax) GetValue() float64 {
	return t.Value.TotalAmount.GetValue()
}
