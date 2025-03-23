package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/invoice"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// ToMHInvoice convierte una ElectronicInvoice a la estructura requerida por Hacienda
func ToMHInvoice(doc *invoice_models.ElectronicInvoice) (*structs.InvoiceDTEResponse, error) {

	dte := &structs.InvoiceDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(doc.Identification),
		Emisor:          common.MapCommonResponseIssuer(doc.Issuer),
		Receptor:        invoice.MapInvoiceResponseReceiver(doc.Receiver),
		Resumen:         invoice.MapInvoiceResponseSummary(doc.InvoiceSummary),
		CuerpoDocumento: invoice.MapInvoiceResponseItem(doc.InvoiceItems),
		Extension:       common.MapCommonResponseExtension(doc.Extension),
	}

	if len(doc.GetRelatedDocuments()) > 0 {
		dte.DocumentoRelacionado = common.MapCommonResponseRelatedDocuments(doc.GetRelatedDocuments())
	}

	if len(doc.GetOtherDocuments()) > 0 {
		dte.OtrosDocumentos = common.MapCommonResponseOtherDocuments(doc.GetOtherDocuments())
	}

	if doc.GetThirdPartySale() != nil {
		dte.VentaTercero = common.MapCommonResponseThirdPartySale(doc.GetThirdPartySale())
	}

	if doc.Appendix != nil {
		dte.Apendice = common.MapCommonResponseAppendix(doc.Appendix)
	}

	return dte, nil
}
