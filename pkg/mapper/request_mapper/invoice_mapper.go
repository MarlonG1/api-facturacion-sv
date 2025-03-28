package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/invoice"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type InvoiceMapper struct{}

func NewInvoiceMapper() *InvoiceMapper {
	return &InvoiceMapper{}
}

// MapToInvoiceData convierte una solicitud de invoice a datos de invoice_models.
func (m *InvoiceMapper) MapToInvoiceData(req *structs.CreateInvoiceRequest, client *dte.IssuerDTE) (*invoice_models.InvoiceData, error) {
	if err := validateInvoiceRequest(req); err != nil {
		return nil, err
	}

	items, err := invoice.MapInvoiceItems(req.Items)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "Error mapping items", err)
	}

	receiver, err := common.MapCommonRequestReceiver(req.Receiver)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "Error mapping receiver", err)
	}

	identification, err := common.MapCommonRequestIdentification(constants.ModeloFacturacionPrevio, 1, constants.FacturaElectronica)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "Error mapping identification", err)
	}

	summary, err := invoice.MapInvoiceRequestSummary(req.Summary)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "Error mapping summary", err)
	}

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "Error mapping issuer", err)
	}

	result := &invoice_models.InvoiceData{
		InputDataCommon: &models.InputDataCommon{
			Issuer:         issuer,
			Identification: identification,
			Receiver:       receiver,
		},
		Items:          items,
		InvoiceSummary: summary,
	}

	if err = mapOptionalFields(req, result); err != nil {
		return nil, err
	}

	return result, nil
}

// validateInvoiceRequest valida los campos requeridos en la solicitud de invoice_models.
func validateInvoiceRequest(req *structs.CreateInvoiceRequest) error {
	if req == nil {
		return shared_error.NewGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "The request is empty", nil)
	}
	if req.Items == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Items")
	}
	if req.Summary == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Summary")
	}

	if req.Receiver != nil {
		if req.Receiver.DocumentType != nil && req.Receiver.DocumentNumber == nil || req.Receiver.DocumentType == nil && req.Receiver.DocumentNumber != nil {
			return dte_errors.NewValidationError("InvalidField", "DocumentType, DocumentNumber. If DocumentType is present, DocumentNumber must be present and vice versa, there fields")
		}
	}
	return nil
}

// mapOptionalFields mapea campos opcionales como ThirdPartySale, Extension, Payments, etc.
func mapOptionalFields(req *structs.CreateInvoiceRequest, result *invoice_models.InvoiceData) error {
	if req.ThirdPartySale != nil {
		thirdPartySale, err := common.MapCommonRequestThirdPartySale(req.ThirdPartySale)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestThirdPartySale", "MapToInvoiceData", "Error mapping third party sale", err)
		}
		result.ThirdPartySale = thirdPartySale
	}

	if req.Extension != nil {
		if req.Extension.VehiculePlate != nil {
			return dte_errors.NewValidationError("InvalidField", "Request->Extension->VehiculePlate")
		}

		extension, err := common.MapCommonRequestExtension(req.Extension)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestExtension", "MapToInvoiceData", "Error mapping extension", err)
		}
		result.Extension = extension
	}

	if req.Payments != nil {
		payments, err := common.MapCommonRequestPaymentsType(req.Payments)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestPaymentsType", "MapToInvoiceData", "Error mapping payments", err)
		}
		result.InvoiceSummary.PaymentTypes = payments
	}

	if req.OtherDocs != nil {
		otherDocs, err := common.MapCommonRequestOtherDocuments(req.OtherDocs)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestOtherDocuments", "MapToInvoiceData", "Error mapping other documents", err)
		}
		result.OtherDocs = otherDocs
	}

	if req.RelatedDocs != nil {
		relatedDocs, err := common.MapCommonRequestRelatedDocuments(req.RelatedDocs)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestRelatedDocuments", "MapToInvoiceData", "Error mapping related documents", err)
		}
		result.RelatedDocs = relatedDocs
	}

	if req.Appendixes != nil {
		appendixes, err := common.MapCommonRequestAppendix(req.Appendixes)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapAppendixes", "MapToInvoiceData", "Error mapping appendixes", err)
		}
		result.Appendixes = appendixes
	}

	return nil
}
