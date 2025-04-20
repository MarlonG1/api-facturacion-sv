package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/credit_note"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type CreditNoteMapper struct{}

func NewCreditNoteMapper() *CreditNoteMapper {
	return &CreditNoteMapper{}
}

// MapToCreditNoteData convierte una solicitud de Nota de Crédito a datos de modelo de dominio.
func (m *CreditNoteMapper) MapToCreditNoteData(req *structs.CreateCreditNoteRequest, client *dte.IssuerDTE) (*credit_note_models.CreditNoteInput, error) {
	if err := validateCreditNoteRequest(req); err != nil {
		return nil, err
	}

	items, err := credit_note.MapCreditNoteItems(req.Items)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("CreditNoteMapper", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Items")
	}

	receiver, err := credit_note.MapCreditNoteRequestReceiver(req.Receiver)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("CreditNoteMapper", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Receiver")
	}

	identification, err := common.MapCommonRequestIdentification(constants.ModeloFacturacionPrevio, 3, constants.NotaCreditoElectronica)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("CreditNoteMapper", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Identification")
	}

	summary, err := credit_note.MapCreditNoteRequestSummary(req.Summary)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("CreditNoteMapper", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Summary")
	}

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("CreditNoteMapper", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Issuer")
	}

	// En notas de crédito, los documentos relacionados son obligatorios y ya validados
	relatedDocs, err := common.MapCommonRequestRelatedDocuments(req.RelatedDocs)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestRelatedDocuments", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->RelatedDocs")
	}

	result := &credit_note_models.CreditNoteInput{
		InputDataCommon: &models.InputDataCommon{
			Issuer:         issuer,
			Identification: identification,
			Receiver:       receiver,
			RelatedDocs:    relatedDocs,
		},
		Items:         items,
		CreditSummary: summary,
	}

	if err = mapCreditNoteOptionalFields(req, result); err != nil {
		return nil, err
	}

	return result, nil
}

// validateCreditNoteRequest valida que la solicitud de Nota de Crédito sea correcta.
func validateCreditNoteRequest(req *structs.CreateCreditNoteRequest) error {
	if req == nil {
		return dte_errors.NewValidationError("RequiredField", "Request")
	}
	if req.Items == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Items")
	}
	if req.Summary == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Summary")
	}
	if req.Receiver == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Receiver")
	}
	// Para notas de crédito, los documentos relacionados son obligatorios
	if req.RelatedDocs == nil || len(req.RelatedDocs) == 0 {
		return dte_errors.NewValidationError("RequiredField", "Request->RelatedDocs")
	}

	for _, doc := range req.RelatedDocs {
		if doc.DocumentType == "" {
			return dte_errors.NewValidationError("RequiredField", "Request->RelatedDocs->DocumentType")
		}
		if doc.DocumentNumber == "" {
			return dte_errors.NewValidationError("RequiredField", "Request->RelatedDocs->DocumentNumber")
		}

		if doc.GenerationType == 0 {
			return dte_errors.NewValidationError("RequiredField", "Request->RelatedDocs->GenerationType")
		}

		if doc.GenerationType == constants.PhysicalDocument && doc.EmissionDate == "" {
			return dte_errors.NewValidationError("InvalidEmissionDateForPhysicalDocument", doc.EmissionDate)
		}
	}

	return nil
}

// mapCreditNoteOptionalFields mapea los campos opcionales de la solicitud de Nota de Crédito.
func mapCreditNoteOptionalFields(req *structs.CreateCreditNoteRequest, result *credit_note_models.CreditNoteInput) error {
	if req.ThirdPartySale != nil {
		thirdPartySale, err := common.MapCommonRequestThirdPartySale(req.ThirdPartySale)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestThirdPartySale", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->ThirdPartySales")
		}
		result.ThirdPartySale = thirdPartySale
	}

	if req.Extension != nil {
		extension, err := common.MapCommonRequestExtension(req.Extension)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestExtension", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Extension")
		}
		result.Extension = extension
	}

	if req.Payments != nil {
		payments, err := common.MapCommonRequestPaymentsType(req.Payments)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestPaymentsType", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->PaymentTypes")
		}
		result.CreditSummary.PaymentTypes = payments
	}

	if req.OtherDocs != nil {
		otherDocs, err := common.MapCommonRequestOtherDocuments(req.OtherDocs)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestOtherDocuments", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->OtherDocs")
		}
		result.OtherDocs = otherDocs
	}

	if req.Appendixes != nil {
		appendixes, err := common.MapCommonRequestAppendix(req.Appendixes)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapAppendixes", "MapToCreditNoteData", err, "ErrorMapping", "CreditNote->Appendixes")
		}
		result.Appendixes = appendixes
	}

	return nil
}
