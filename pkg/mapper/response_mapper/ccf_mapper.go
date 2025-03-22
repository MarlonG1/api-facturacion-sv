package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/ccf"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToMHCreditFiscalInvoice(doc *ccf_models.CreditFiscalDocument) (*structs.CCFDTEResponse, error) {
	dte := &structs.CCFDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(doc.Identification),
		Emisor:          common.MapCommonResponseIssuer(doc.Issuer),
		Receptor:        common.MapCommonResponseReceiver(doc.Receiver),
		Resumen:         ccf.MapCCFResponseSummary(doc.CreditSummary),
		CuerpoDocumento: ccf.MapCCFResponseItem(doc.CreditItems),
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

	if doc.GetAppendix() != nil {
		dte.Apendice = common.MapCommonResponseAppendix(doc.GetAppendix())
	}

	return dte, nil
}
