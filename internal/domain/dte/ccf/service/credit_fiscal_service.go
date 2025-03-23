package service

import (
	"context"
	"fmt"
	localInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	buisnessValidator "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type creditFiscalService struct {
	validator     *validator.CCFRulesValidator
	seqNumberRepo ports.SequentialNumberRepositoryPort
}

// NewCCFService Crea un nuevo servicio Comprobante de Crédito Fiscal.
func NewCCFService(seqNumberRepo ports.SequentialNumberRepositoryPort) localInterfaces.CCFManager {
	return &creditFiscalService{
		validator:     validator.NewCCFRulesValidator(nil),
		seqNumberRepo: seqNumberRepo,
	}
}

func (c *creditFiscalService) Create(ctx context.Context, data *ccf_models.CCFData, branchID uint) (*ccf_models.CreditFiscalDocument, error) {
	baseDoc := createBaseDocument(data)

	creditFiscalDocument := &ccf_models.CreditFiscalDocument{
		DTEDocument:   baseDoc,
		CreditItems:   data.Items,
		CreditSummary: *data.CreditSummary,
	}

	if err := c.Validate(creditFiscalDocument); err != nil {
		logs.Error("Failed to validate credit fiscal document basic validation", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	if err := buisnessValidator.ValidateDTEDocument(creditFiscalDocument); err != nil {
		logs.Error("Failed to validate credit fiscal document generic validations", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	if err := c.generateCodeAndIdentifiers(ctx, creditFiscalDocument, branchID); err != nil {
		return nil, err
	}

	return creditFiscalDocument, nil
}

func (c *creditFiscalService) Validate(ccf *ccf_models.CreditFiscalDocument) error {
	c.validator = validator.NewCCFRulesValidator(ccf)
	err := c.validator.Validate()
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"CCFService",
			"Validate",
			"validation failed, check the error for more details",
			err,
		)
	}
	return nil
}

func (c *creditFiscalService) generateControlNumber(ctx context.Context, ccf *ccf_models.CreditFiscalDocument, branchID uint) error {
	establishmentCode := ccf.Issuer.GetEstablishmentCode()
	posCode := ccf.Issuer.GetPOSCode()
	defaultCode := "0000"

	if posCode == nil {
		posCode = &defaultCode
	}
	if establishmentCode == nil {
		establishmentCode = &defaultCode
	}

	correlativeNumber, err := c.seqNumberRepo.GetNext(
		ctx,
		constants.CCFElectronico,
		branchID,
	)
	if err != nil {
		return err
	}

	controlNumber := fmt.Sprintf("DTE-%s-%s%s-%015d",
		constants.CCFElectronico,
		*establishmentCode,
		*posCode,
		correlativeNumber,
	)

	err = ccf.Identification.SetControlNumber(controlNumber)
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"InvoiceService",
			"GenerateControlNumber",
			"failed to set control number",
			err,
		)
	}
	return nil
}

func (c *creditFiscalService) generateCodeAndIdentifiers(ctx context.Context, ccf *ccf_models.CreditFiscalDocument, branchID uint) error {
	if err := c.generateControlNumber(ctx, ccf, branchID); err != nil {
		return err
	}
	return ccf.Identification.GenerateCode()
}

func (c *creditFiscalService) IsValid(ccf *ccf_models.CreditFiscalDocument) bool {
	return c.Validate(ccf) == nil
}

func createBaseDocument(data *ccf_models.CCFData) *models.DTEDocument {
	var extInterface interfaces.Extension
	var thirdPartySale interfaces.ThirdPartySale
	var appendixes []interfaces.Appendix
	var relatedDocuments []interfaces.RelatedDocument
	var otherDocuments []interfaces.OtherDocuments
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
