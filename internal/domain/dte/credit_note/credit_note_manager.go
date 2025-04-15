package credit_note

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	creditValidator "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/validator"
	seqPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type creditNoteService struct {
	validator        *creditValidator.CreditNoteRulesValidator
	seqNumberManager seqPorts.SequentialNumberManager
	dteManager       seqPorts.DTEManager
}

// NewCreditNoteService Crea un nuevo servicio de Nota de Crédito.
func NewCreditNoteService(seqNumberManager seqPorts.SequentialNumberManager, dteManager seqPorts.DTEManager) CreditNoteManager {
	return &creditNoteService{
		validator:        creditValidator.NewCreditNoteRulesValidator(nil),
		seqNumberManager: seqNumberManager,
		dteManager:       dteManager,
	}
}

// Create Crea una nueva Nota de Crédito electrónica con base en los datos proporcionados.
func (s *creditNoteService) Create(ctx context.Context, data *credit_note_models.CreditNoteInput, branchID uint) (*credit_note_models.CreditNoteModel, error) {
	// 1. Validar la existencia de documentos relacionados
	if err := s.validateRelatedDocs(ctx, data, branchID); err != nil {
		logs.Error("Failed to validate related documents", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// 2. Crear el documento base
	baseDoc := createBaseDocument(data)
	creditNote := &credit_note_models.CreditNoteModel{
		DTEDocument:   baseDoc,
		CreditItems:   data.Items,
		CreditSummary: *data.CreditSummary,
	}

	// 3. Validar el documento base
	if err := s.Validate(creditNote); err != nil {
		logs.Error("Failed to validate credit note document basic validation", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// 4. Validar contra reglas principales de negocio
	if err := validator.ValidateDTEDocument(creditNote); err != nil {
		logs.Error("Failed to validate credit note document generic validations", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// 5. Validar documento totales
	for _, doc := range creditNote.RelatedDocuments {
		if doc.GetGenerationType() == constants.ElectronicDocument {
			if err := s.dteManager.ValidateForCreditNote(ctx, branchID, doc.GetDocumentNumber(), creditNote); err != nil {
				logs.Error("Failed to validate credit note document totals", map[string]interface{}{"error": err.Error()})
				return nil, err
			}
		}
	}

	// 6. Generar el número de control y el código UUID
	if err := s.generateCodeAndIdentifiers(ctx, creditNote, branchID); err != nil {
		return nil, err
	}

	return creditNote, nil
}

// validateRelatedDocs verifica que los documentos relacionados existan en la base de datos
func (s *creditNoteService) validateRelatedDocs(ctx context.Context, data *credit_note_models.CreditNoteInput, branchID uint) error {
	if data.RelatedDocs == nil || len(data.RelatedDocs) == 0 {
		return shared_error.NewGeneralServiceError(
			"CreditNoteService",
			"validateRelatedDocs",
			"no related documents provided",
			nil,
		)
	}

	// Verificar que cada documento relacionado exista en la base de datos
	for _, relatedDoc := range data.RelatedDocs {
		_, err := s.dteManager.GetByGenerationCode(ctx, branchID, relatedDoc.GetDocumentNumber())
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate Valida una Nota de Crédito electrónica con base en las reglas de negocio.
func (s *creditNoteService) Validate(creditNote *credit_note_models.CreditNoteModel) error {
	s.validator = creditValidator.NewCreditNoteRulesValidator(creditNote)
	err := s.validator.Validate()
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"CreditNoteService",
			"Validate",
			"validation failed, check the error for more details",
			err,
		)
	}
	return nil
}

// IsValid Comprueba si una Nota de Crédito electrónica es válida según las reglas de negocio.
func (s *creditNoteService) IsValid(creditNote *credit_note_models.CreditNoteModel) bool {
	return s.Validate(creditNote) == nil
}

// generateControlNumber Genera un número de control único para la Nota de Crédito.
func (s *creditNoteService) generateControlNumber(ctx context.Context, creditNote *credit_note_models.CreditNoteModel, branchID uint) error {
	establishmentCode := creditNote.Issuer.GetEstablishmentCode()
	posCode := creditNote.Issuer.GetPOSCode()

	controlNumber, err := s.seqNumberManager.GetNextControlNumber(
		ctx,
		constants.NotaCreditoElectronica,
		branchID,
		posCode,
		establishmentCode,
	)
	if err != nil {
		return err
	}

	err = creditNote.Identification.SetControlNumber(controlNumber)
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"CreditNoteService",
			"GenerateControlNumber",
			"failed to set control number",
			err,
		)
	}
	return nil
}

// generateCodeAndIdentifiers Genera el código UUID y número de control de la Nota de Crédito.
func (s *creditNoteService) generateCodeAndIdentifiers(ctx context.Context, creditNote *credit_note_models.CreditNoteModel, branchID uint) error {
	err := creditNote.Identification.GenerateCode()
	if err != nil {
		return err
	}

	return s.generateControlNumber(ctx, creditNote, branchID)
}

// createBaseDocument Crea un documento base para la Nota de Crédito electrónica.
func createBaseDocument(data *credit_note_models.CreditNoteInput) *models.DTEDocument {
	var extInterface interfaces.Extension
	var thirdPartySale interfaces.ThirdPartySale
	var appendixes []interfaces.Appendix
	var otherDocuments []interfaces.OtherDocuments
	var relatedDocuments []interfaces.RelatedDocument

	baseItems := make([]interfaces.Item, len(data.Items))
	for i, item := range data.Items {
		baseItems[i] = &item
	}

	if data.Appendixes != nil {
		for _, appendix := range data.Appendixes {
			appendixes = append(appendixes, &appendix)
		}
	}

	if data.Extension != nil {
		extInterface = data.Extension
	}

	if data.RelatedDocs != nil {
		for _, relatedDoc := range data.RelatedDocs {
			relatedDocuments = append(relatedDocuments, &relatedDoc)
		}
	}

	if data.OtherDocs != nil {
		for _, otherDoc := range data.OtherDocs {
			otherDocuments = append(otherDocuments, &otherDoc)
		}
	}

	if data.ThirdPartySale != nil {
		thirdPartySale = data.ThirdPartySale
	}

	return &models.DTEDocument{
		Identification:   data.Identification,
		Issuer:           data.Issuer,
		Receiver:         data.Receiver,
		Items:            baseItems,
		RelatedDocuments: relatedDocuments,
		OtherDocuments:   otherDocuments,
		Summary:          data.CreditSummary.Summary,
		ThirdPartySale:   thirdPartySale,
		Extension:        extInterface,
		Appendix:         appendixes,
	}
}
