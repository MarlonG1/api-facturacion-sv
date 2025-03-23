package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// PaymentType es una estructura que representa un tipo de pago de un DTE, contiene Code, Amount, Reference, Term y Period
type PaymentType struct {
	Code      financial.PaymentType  `json:"code"`
	Amount    financial.Amount       `json:"amount"`
	Reference string                 `json:"reference"`
	Term      *financial.PaymentTerm `json:"term,omitempty"`
	Period    *int                   `json:"period,omitempty"`
}

func (p *PaymentType) GetCode() string {
	return p.Code.GetValue()
}
func (p *PaymentType) GetAmount() float64 {
	return p.Amount.GetValue()
}
func (p *PaymentType) GetReference() string {
	return p.Reference
}
func (p *PaymentType) GetTerm() *string {
	if p.Term == nil {
		return nil
	}

	return utils.ToStringPointer(p.Term.GetValue())
}
func (p *PaymentType) GetPeriod() *int {
	if p == nil {
		return nil
	}

	return p.Period
}
func (p *PaymentType) GetPeriodPointer() *int {
	return p.Period
}
