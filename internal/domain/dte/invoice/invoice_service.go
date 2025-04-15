package invoice

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	buisnessValidator "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type invoiceService struct {
	validator        *validator.InvoiceRulesValidator
	seqNumberManager dte_documents.SequentialNumberManager
}

// NewInvoiceService Crea un nuevo servicio de facturas electrónicas.
func NewInvoiceService(seqNumberManager dte_documents.SequentialNumberManager) ports.DTEService {
	return &invoiceService{
		validator:        validator.NewInvoiceRulesValidator(nil),
		seqNumberManager: seqNumberManager,
	}
}

// Create Crea una nueva invoice electrónica con base en los datos proporcionados.
func (s *invoiceService) Create(ctx context.Context, input interface{}, branchID uint) (interface{}, error) {
	data := input.(*invoice_models.InvoiceData)
	baseDoc := createBaseDocument(data)

	invoice := &invoice_models.ElectronicInvoice{
		DTEDocument:    baseDoc,
		InvoiceItems:   data.Items,
		InvoiceSummary: *data.InvoiceSummary,
	}

	if err := s.validate(invoice); err != nil {
		return nil, err
	}

	if err := buisnessValidator.ValidateDTEDocument(invoice); err != nil {
		return nil, err
	}

	if err := s.generateCodeAndIdentifiers(ctx, invoice, branchID); err != nil {
		return nil, err
	}

	return invoice, nil
}

// Validate Valida una invoice electrónica con base en las reglas de negocio.
func (s *invoiceService) validate(invoice *invoice_models.ElectronicInvoice) error {
	s.validator = validator.NewInvoiceRulesValidator(invoice)
	err := s.validator.Validate()
	if err != nil {
		return shared_error.NewGeneralServiceError(
			"InvoiceService",
			"Validate",
			"validation failed, check the error for more details",
			err,
		)
	}
	return nil
}

// createBaseDocument Crea un documento base para la invoice electrónica.
func createBaseDocument(data *invoice_models.InvoiceData) *models.DTEDocument {
	var extInterface interfaces.Extension
	var appendixes []interfaces.Appendix
	var relatedDocuments []interfaces.RelatedDocument
	var otherDocuments []interfaces.OtherDocuments
	var thirdPartySale interfaces.ThirdPartySale
	receiver := &models.Receiver{
		Address: &models.Address{},
	}

	baseItems := make([]interfaces.Item, len(data.Items))
	for i, item := range data.Items {
		baseItems[i] = &item
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
		Receiver:         receiver,
		Items:            baseItems,
		Summary:          data.InvoiceSummary.Summary,
		Extension:        extInterface,
		RelatedDocuments: relatedDocuments,
		OtherDocuments:   otherDocuments,
		ThirdPartySale:   thirdPartySale,
		Appendix:         appendixes,
	}
}

// generateControlNumber Genera un número de control único para la invoice.
func (s *invoiceService) generateControlNumber(ctx context.Context, invoice *invoice_models.ElectronicInvoice, branchID uint) error {
	establishmentCode := invoice.Issuer.GetEstablishmentCode()
	posCode := invoice.Issuer.GetPOSCode()

	controlNumber, err := s.seqNumberManager.GetNextControlNumber(
		ctx,
		constants.FacturaElectronica,
		branchID,
		posCode,
		establishmentCode,
	)
	if err != nil {
		return err
	}

	err = invoice.Identification.SetControlNumber(controlNumber)
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

// generateCodeAndIdentifiers Genera el código UUID y número de control de la invoice.
func (s *invoiceService) generateCodeAndIdentifiers(ctx context.Context, invoice *invoice_models.ElectronicInvoice, branchID uint) error {
	if err := s.generateControlNumber(ctx, invoice, branchID); err != nil {
		return err
	}
	return invoice.Identification.GenerateCode()
}
