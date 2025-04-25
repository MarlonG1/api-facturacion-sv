package retention

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type retentionService struct {
	validator        *validator.RetentionRulesValidator
	dteManager       dte_documents.DTEManager
	seqNumberManager dte_documents.SequentialNumberManager
}

// NewRetentionService crea una nueva instancia de RetentionManager
func NewRetentionService(seqNumberManager dte_documents.SequentialNumberManager, dteManager dte_documents.DTEManager) ports.DTEService {
	return &retentionService{
		validator:        validator.NewRetentionRulesValidator(nil),
		seqNumberManager: seqNumberManager,
		dteManager:       dteManager,
	}
}

func (s *retentionService) Create(ctx context.Context, input interface{}, branchID uint) (interface{}, error) {
	data := input.(*retention_models.InputRetentionData)

	// 1. Crear el documento base para la retención
	data.RetentionSummary.TotalIVARetentionLetters = utils.InLetters(data.RetentionSummary.TotalIVARetention.GetValue())
	baseDoc := createBaseDocument(data)
	retention := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   data.RetentionItems,
		RetentionSummary: data.RetentionSummary,
	}

	// 2. Validar el documento de retention generado
	err := s.validate(retention)
	if err != nil {
		return nil, err
	}

	// 3. Generar el codigo de generacion y el numero de control
	if err := s.generateCodeAndIdentifiers(ctx, retention, branchID); err != nil {
		return nil, err
	}

	return retention, nil

}

func (s *retentionService) validate(retention *retention_models.RetentionModel) error {
	s.validator = validator.NewRetentionRulesValidator(retention)
	err := s.validator.Validate()
	if err != nil {
		return shared_error.NewFormattedGeneralServiceWithError(
			"RetentionService",
			"Validate",
			err,
			"ValidationFailed",
		)
	}

	return nil
}

// createBaseDocument Crea un documento base para la invoice electrónica.
func createBaseDocument(data *retention_models.InputRetentionData) *models.DTEDocument {
	var extInterface interfaces.Extension
	var appendixes []interfaces.Appendix
	var items []interfaces.Item
	var relatedDocuments []interfaces.RelatedDocument
	var otherDocuments []interfaces.OtherDocuments
	var thirdPartySale interfaces.ThirdPartySale
	receiver := &models.Receiver{
		Address: &models.Address{},
	}

	if data.Appendixes != nil {
		for _, appendix := range data.Appendixes {
			appendixes = append(appendixes, &appendix)
		}
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

	if data.Extension != nil {
		extInterface = data.Extension
	}

	if data.ThirdPartySale != nil {
		thirdPartySale = data.ThirdPartySale
	}

	if data.Receiver != nil {
		receiver = data.Receiver
	}

	return &models.DTEDocument{
		Identification:   data.Identification,
		Issuer:           data.Issuer,
		Items:            items,
		Receiver:         receiver,
		Extension:        extInterface,
		RelatedDocuments: relatedDocuments,
		OtherDocuments:   otherDocuments,
		ThirdPartySale:   thirdPartySale,
		Appendix:         appendixes,
	}
}

func (s *retentionService) generateCodeAndIdentifiers(ctx context.Context, retention *retention_models.RetentionModel, branchID uint) error {
	if err := s.generateControlNumber(ctx, retention, branchID); err != nil {
		return err
	}
	return retention.Identification.GenerateCode()
}

// generateControlNumber Genera un número de control único para la invoice.
func (s *retentionService) generateControlNumber(ctx context.Context, retention *retention_models.RetentionModel, branchID uint) error {
	establishmentCode := retention.Issuer.GetEstablishmentCode()
	posCode := retention.Issuer.GetPOSCode()

	controlNumber, err := s.seqNumberManager.GetNextControlNumber(
		ctx,
		constants.ComprobanteRetencionElectronico,
		branchID,
		posCode,
		establishmentCode,
	)
	if err != nil {
		return err
	}

	err = retention.Identification.SetControlNumber(controlNumber)
	if err != nil {
		return shared_error.NewFormattedGeneralServiceWithError(
			"RetentionManager",
			"GenerateControlNumber",
			err,
			"FailedToSetControlNumber",
		)
	}
	return nil
}
