package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/credit_note"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToMHCreditNote(doc *credit_note_models.CreditNoteModel) (*structs.CreditNoteDTEResponse, error) {
	dte := &structs.CreditNoteDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(doc.Identification),
		Receptor:        common.MapCommonResponseReceiver(doc.Receiver),
		Emisor:          credit_note.MapCreditNoteIssuer(doc.Issuer),
		Resumen:         credit_note.MapCreditNoteResponseSummary(doc.CreditSummary),
		CuerpoDocumento: credit_note.MapCreditNoteResponseItem(doc.CreditItems),
		Extension:       credit_note.MapCreditNoteResponseExtension(doc.Extension),
	}

	// En Nota de Crédito, los documentos relacionados siempre deben existir
	dte.DocumentoRelacionado = common.MapCommonResponseRelatedDocuments(doc.GetRelatedDocuments())

	if doc.GetThirdPartySale() != nil {
		dte.VentaTercero = common.MapCommonResponseThirdPartySale(doc.GetThirdPartySale())
	}

	if doc.GetAppendix() != nil {
		dte.Apendice = common.MapCommonResponseAppendix(doc.GetAppendix())
	}

	return dte, nil
}
