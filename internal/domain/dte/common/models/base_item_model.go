package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
)

// Item es una estructura que representa un item de un DTE, contiene Number, Type, Code, Description, Quantity,
// UnitMeasure, UnitPrice, Discount, Taxes, TaxCode y RelatedDoc
type Item struct {
	Number      item.ItemNumber    `json:"number"`
	Type        item.ItemType      `json:"type"`
	Description string             `json:"description"`
	Quantity    item.Quantity      `json:"quantity"`
	UnitMeasure item.UnitMeasure   `json:"unitMeasure"`
	UnitPrice   financial.Amount   `json:"unitPrice"`
	Discount    financial.Discount `json:"discount"`
	Taxes       []string           `json:"taxes"`
	Code        *item.ItemCode     `json:"code,omitempty"`
	TaxCode     *financial.TaxType `json:"taxCode,omitempty"`
	RelatedDoc  *string            `json:"relatedDoc,omitempty"`
}

func (i *Item) GetUnitMeasure() int {
	return i.UnitMeasure.GetValue()
}
func (i *Item) GetNumber() int {
	return i.Number.GetValue()
}
func (i *Item) GetQuantity() float64 {
	return i.Quantity.GetValue()
}
func (i *Item) GetItemCode() string {
	return i.Code.GetValue()
}
func (i *Item) GetDescription() string {
	return i.Description
}
func (i *Item) GetType() int {
	return i.Type.GetValue()
}
func (i *Item) GetUnitPrice() float64 {
	return i.UnitPrice.GetValue()
}
func (i *Item) GetDiscount() float64 {
	return i.Discount.GetValue()
}

func (i *Item) GetTaxes() []string {
	if len(i.Taxes) == 0 {
		return nil
	}
	return i.Taxes
}
func (i *Item) GetRelatedDoc() *string {
	return i.RelatedDoc
}
