package service

import (
	"context"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	buisnessValidator "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator"
	localInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/validator"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type invoiceService struct {
	validator     *validator.InvoiceRulesValidator
	seqNumberRepo ports.SequentialNumberRepositoryPort
}

// NewInvoiceService Crea un nuevo servicio de facturas electrónicas.
func NewInvoiceService(seqNumberRepo ports.SequentialNumberRepositoryPort) localInterfaces.InvoiceManager {
	return &invoiceService{
		validator:     validator.NewInvoiceRulesValidator(nil),
		seqNumberRepo: seqNumberRepo,
	}
}

// Create Crea una nueva invoice electrónica con base en los datos proporcionados.
func (s *invoiceService) Create(ctx context.Context, data *invoice_models.InvoiceData, branchID uint) (*invoice_models.ElectronicInvoice, error) {
	baseDoc := createBaseDocument(data)

	invoice := &invoice_models.ElectronicInvoice{
		DTEDocument:    baseDoc,
		InvoiceItems:   data.Items,
		InvoiceSummary: *data.InvoiceSummary,
	}

	if err := s.Validate(invoice); err != nil {
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
func (s *invoiceService) Validate(invoice *invoice_models.ElectronicInvoice) error {
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

// IsValid Comprueba si una invoice electrónica es válida según las reglas de negocio.
func (s *invoiceService) IsValid(invoice *invoice_models.ElectronicInvoice) bool {
	return s.Validate(invoice) == nil
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
	defaultCode := "0000"

	if posCode == nil {
		posCode = &defaultCode
	}
	if establishmentCode == nil {
		establishmentCode = &defaultCode
	}

	correlativeNumber, err := s.seqNumberRepo.GetNext(
		ctx,
		constants.FacturaElectronica,
		branchID,
	)
	if err != nil {
		return err
	}

	controlNumber := fmt.Sprintf("DTE-%s-%s%s-%015d",
		constants.FacturaElectronica,
		*establishmentCode,
		*posCode,
		correlativeNumber,
	)

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
