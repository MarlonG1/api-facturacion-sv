package service

import (
	"context"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	seqPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	retentionInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/validator"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type RetentionManager struct {
	validator        *validator.RetentionRulesValidator
	dteManager       seqPorts.DTEManager
	seqNumberManager seqPorts.SequentialNumberManager
}

// NewRetentionManager crea una nueva instancia de RetentionManager
func NewRetentionManager(seqNumberManager seqPorts.SequentialNumberManager, dteManager seqPorts.DTEManager) retentionInterfaces.RetentionManager {
	return &RetentionManager{
		validator:        validator.NewRetentionRulesValidator(nil),
		seqNumberManager: seqNumberManager,
		dteManager:       dteManager,
	}
}

func (s *RetentionManager) Create(ctx context.Context, data *retention_models.InputRetentionData, branchID uint, isAllPhysical bool) (*retention_models.RetentionModel, error) {
	// 1. Verificar si todos los documentos son físicos, si no lo son, obtener los detalles de cada documento
	if !isAllPhysical {
		for i := range data.RetentionItems {

			// 1.1 Obtener el documento DTE correspondiente al número de documento
			dte, err := s.dteManager.GetByGenerationCode(ctx, branchID, data.RetentionItems[i].DocumentNumber.GetValue())
			if err != nil {
				return nil, err
			}

			// 1.2 Verificar si el tipo de DTE es válido para retención
			if !constants.ValidRetentionDTETypes[dte.Details.DTEType] {
				return nil, shared_error.NewGeneralServiceError("RetentionUseCase", "Create", fmt.Sprintf("dte type for %s is not valid for retention", data.RetentionItems[i].DocumentNumber.GetValue()), nil)
			}

			// 1.3 Verificar si el documento tiene detalles
			err = s.extractSummaryData(&data.RetentionItems[i], dte)
			if err != nil {
				return nil, err
			}

		}

		// 1.4 Calcular el resumen de la retención
		s.calculateSummary(data)
	}

	// 2. Crear el documento base para la retención
	data.RetentionSummary.TotalIVARetentionLetters = utils.InLetters(data.RetentionSummary.TotalIVARetention.GetValue())
	baseDoc := createBaseDocument(data)
	retention := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   data.RetentionItems,
		RetentionSummary: data.RetentionSummary,
	}

	// 3. Validar el documento de retention generado
	err := s.Validate(retention)
	if err != nil {
		return nil, err
	}

	// 4. Generar el codigo de generacion y el numero de control
	if err := s.generateCodeAndIdentifiers(ctx, retention, branchID); err != nil {
		return nil, err
	}

	return retention, nil

}

func (s *RetentionManager) Validate(retention *retention_models.RetentionModel) error {
	s.validator = validator.NewRetentionRulesValidator(retention)
	err := s.validator.Validate()
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"RetentionService",
			"Validate",
			"validation failed, check the error for more details",
			err,
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

func (s *RetentionManager) IsValid(retention *retention_models.RetentionModel) bool {
	return s.Validate(retention) == nil
}

func (s *RetentionManager) calculateSummary(retention *retention_models.InputRetentionData) {
	var totalIva, totalAmount financial.Amount

	for _, item := range retention.RetentionItems {
		totalIva.Add(&item.RetentionIVA)
		totalAmount.Add(&item.RetentionAmount)
	}

	retention.RetentionSummary.TotalIVARetention = totalIva
	retention.RetentionSummary.TotalSubjectRetention = totalAmount
}

func (s *RetentionManager) extractSummaryData(item *retention_models.RetentionItem, doc *dte.DTEDocument) error {
	extractor, err := utils.ExtractAuxiliarSummaryFromStringJSON(doc.Details.JSONData)
	if err != nil {
		return err
	}

	item.RetentionAmount = *financial.NewValidatedAmount(extractor.Summary.SubTotal)
	item.RetentionIVA = *financial.NewValidatedAmount(extractor.Summary.IvaRetention)
	item.EmissionDate = *temporal.NewValidatedEmissionDate(doc.CreatedAt)
	item.DTEType = *document.NewValidatedDTEType(doc.Details.DTEType)

	return nil
}

func (s *RetentionManager) generateCodeAndIdentifiers(ctx context.Context, retention *retention_models.RetentionModel, branchID uint) error {
	if err := s.generateControlNumber(ctx, retention, branchID); err != nil {
		return err
	}
	return retention.Identification.GenerateCode()
}

// generateControlNumber Genera un número de control único para la invoice.
func (s *RetentionManager) generateControlNumber(ctx context.Context, retention *retention_models.RetentionModel, branchID uint) error {
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
		return shared_error.NewGeneralServiceError(
			"RetentionManager",
			"GenerateControlNumber",
			"failed to set control number",
			err,
		)
	}
	return nil
}
