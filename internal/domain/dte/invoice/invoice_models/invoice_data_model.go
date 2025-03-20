package invoice_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type InvoiceData struct {
	*models.InputDataCommon
	Items          []InvoiceItem
	InvoiceSummary *InvoiceSummary
}
