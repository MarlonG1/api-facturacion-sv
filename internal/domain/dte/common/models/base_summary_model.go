package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

// Summary es una estructura que representa un resumen de un DTE, contiene TotalNonSubject, TotalExempt, TotalTaxed, SubTotal, NonSubjectDiscount,
// ExemptDiscount, DiscountPercentage, TotalDiscount, TotalTaxes, SubTotalOperation, TotalOperation, TotalNonTaxed, PaymentTypes, OperationCondition y ElectronicPayment
type Summary struct {
	TotalNonSubject    financial.Amount           `json:"totalNonSubject"`
	TotalExempt        financial.Amount           `json:"totalExempt"`
	TotalTaxed         financial.Amount           `json:"totalTaxed"`
	SubTotal           financial.Amount           `json:"subTotal"`
	SubTotalSales      financial.Amount           `json:"subTotalSales"`
	NonSubjectDiscount financial.Amount           `json:"nonSubjectDiscount"`
	ExemptDiscount     financial.Amount           `json:"exemptDiscount"`
	DiscountPercentage financial.Discount         `json:"discountPercentage"`
	TotalDiscount      financial.Amount           `json:"totalDiscount"`
	TotalOperation     financial.Amount           `json:"totalOperation"`
	TotalNonTaxed      financial.Amount           `json:"totalNotTaxed"`
	OperationCondition financial.PaymentCondition `json:"operation_condition"`
	TotalToPay         financial.Amount           `json:"totalToPay"`
	TotalTaxes         []interfaces.Tax           `json:"taxes,omitempty"`
	TotalInWords       *string                    `json:"totalInWords,omitempty"`
	ElectronicPayment  *string                    `json:"electronicPayment,omitempty"`
	PaymentTypes       []interfaces.PaymentType   `json:"payments,omitempty"`
}

func (s *Summary) GetTotalNonSubject() float64 {
	return s.TotalNonSubject.GetValue()
}
func (s *Summary) GetTotalExempt() float64 {
	return s.TotalExempt.GetValue()
}
func (s *Summary) GetTotalTaxed() float64 {
	return s.TotalTaxed.GetValue()
}
func (s *Summary) GetSubTotal() float64 {
	return s.SubTotal.GetValue()
}
func (s *Summary) GetNonSubjectDiscount() float64 {
	return s.NonSubjectDiscount.GetValue()
}
func (s *Summary) GetExemptDiscount() float64 {
	return s.ExemptDiscount.GetValue()
}
func (s *Summary) GetDiscountPercentage() float64 {
	return s.DiscountPercentage.GetValue()
}
func (s *Summary) GetTotalDiscount() float64 {
	return s.TotalDiscount.GetValue()
}
func (s *Summary) GetTotalTaxes() []interfaces.Tax {
	return s.TotalTaxes
}
func (s *Summary) GetTotalOperation() float64 {
	return s.TotalOperation.GetValue()
}
func (s *Summary) GetTotalNotTaxed() float64 {
	return s.TotalNonTaxed.GetValue()
}
func (s *Summary) GetPaymentTypes() []interfaces.PaymentType {
	return s.PaymentTypes
}
func (s *Summary) GetOperationCondition() int {
	return s.OperationCondition.GetValue()
}
func (s *Summary) GetSubtotalSales() float64 {
	return s.SubTotalSales.GetValue()
}
func (s *Summary) GetTotalToPay() float64 {
	return s.TotalToPay.GetValue()
}
func (s *Summary) GetTotalInWords() *string {
	return s.TotalInWords
}
func (s *Summary) GetElectronicPayment() *string {
	return s.ElectronicPayment
}
