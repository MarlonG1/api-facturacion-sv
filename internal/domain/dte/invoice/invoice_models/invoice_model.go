package invoice_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type ElectronicInvoice struct {
	*models.DTEDocument `json:"*Models.DTEDocument"` // herencia de la base
	InvoiceItems        []InvoiceItem                `json:"invoiceItems"`   // cuerpoDocumento
	InvoiceSummary      InvoiceSummary               `json:"invoiceSummary"` // resumen
	State               string                       `json:"state"`          // estado
}
