package invoice_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

// InvoiceItem representa un item de la invoice electrónica de venta
type InvoiceItem struct {
	*models.Item
	NonSubjectSale financial.Amount `json:"nonSubjectSale"` // Venta no sujeta
	ExemptSale     financial.Amount `json:"exemptSale"`     // Venta exenta
	TaxedSale      financial.Amount `json:"taxedSale"`      // Venta gravada
	SuggestedPrice financial.Amount `json:"suggestedPrice"` // Precio de venta sugerido
	NonTaxed       financial.Amount `json:"nonTaxed"`       // Monto no gravado
	IVAItem        financial.Amount `json:"ivaItem"`        // IVA del item
}
