package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
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

func (t *Tax) SetTotalAmount(totalAmount float64) error {
	taObj, err := financial.NewAmount(totalAmount)
	if err != nil {
		return err
	}
	if t.Value == nil {
		t.Value = &TaxAmount{}
	}
	t.Value.TotalAmount = *taObj
	return nil
}

func (t *Tax) SetCode(code string) error {
	codeObj, err := financial.NewTaxType(code)
	if err != nil {
		return err
	}
	t.Code = *codeObj
	return nil
}

func (t *Tax) SetDescription(description string) error {
	if description == "" {
		return dte_errors.NewValidationError("RequiredField", "Description")
	}
	t.Description = description
	return nil
}

func (t *Tax) SetValue(value float64) error {
	return t.SetTotalAmount(value)
}
