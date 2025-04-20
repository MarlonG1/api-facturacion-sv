package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/credit_note"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToMHCreditNote(doc interface{}) *structs.CreditNoteDTEResponse {

	cast := doc.(*credit_note_models.CreditNoteModel)
	dte := &structs.CreditNoteDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(cast.Identification),
		Receptor:        common.MapCommonResponseReceiver(cast.Receiver),
		Emisor:          credit_note.MapCreditNoteIssuer(cast.Issuer),
		Resumen:         credit_note.MapCreditNoteResponseSummary(cast.CreditSummary),
		CuerpoDocumento: credit_note.MapCreditNoteResponseItem(cast.CreditItems),
		Extension:       credit_note.MapCreditNoteResponseExtension(cast.Extension),
	}

	// En Nota de Crédito, los documentos relacionados siempre deben existir
	dte.DocumentoRelacionado = common.MapCommonResponseRelatedDocuments(cast.GetRelatedDocuments())

	if cast.GetThirdPartySale() != nil {
		dte.VentaTercero = common.MapCommonResponseThirdPartySale(cast.GetThirdPartySale())
	}

	if cast.GetAppendix() != nil {
		dte.Apendice = common.MapCommonResponseAppendix(cast.GetAppendix())
	}

	return dte
}
