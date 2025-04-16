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
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvoiceMapper", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Items")
	}

	receiver, err := common.MapCommonRequestReceiver(req.Receiver)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvoiceMapper", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Receiver")
	}

	identification, err := common.MapCommonRequestIdentification(constants.ModeloFacturacionPrevio, 1, constants.FacturaElectronica)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvoiceMapper", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Identification")
	}

	summary, err := invoice.MapInvoiceRequestSummary(req.Summary)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvoiceMapper", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Summary")
	}

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvoiceMapper", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Issuer")
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
		return dte_errors.NewValidationError("RequiredField", "Request")
	}
	if req.Items == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Items")
	}
	if req.Summary == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Summary")
	}

	if req.Receiver != nil {
		if req.Receiver.DocumentType != nil && req.Receiver.DocumentNumber == nil || req.Receiver.DocumentType == nil && req.Receiver.DocumentNumber != nil {
			return shared_error.NewFormattedGeneralServiceError("InvoiceMapper", "MapToInvoiceData", "InvalidDocumentTypeAndNumber")
		}
	}
	return nil
}

// mapOptionalFields mapea campos opcionales como ThirdPartySale, Extension, Payments, etc.
func mapOptionalFields(req *structs.CreateInvoiceRequest, result *invoice_models.InvoiceData) error {
	if req.ThirdPartySale != nil {
		thirdPartySale, err := common.MapCommonRequestThirdPartySale(req.ThirdPartySale)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestThirdPartySale", "MapToInvoiceData", err, "ErrorMapping", "Invoice->ThirdPartySales")
		}
		result.ThirdPartySale = thirdPartySale
	}

	if req.Extension != nil {
		if req.Extension.VehiculePlate != nil {
			return dte_errors.NewValidationError("InvalidField", "Request->Extension->VehiculePlate")
		}

		extension, err := common.MapCommonRequestExtension(req.Extension)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestExtension", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Extension")
		}
		result.Extension = extension
	}

	if req.Payments != nil {
		payments, err := common.MapCommonRequestPaymentsType(req.Payments)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestPaymentsType", "MapToInvoiceData", err, "ErrorMapping", "Invoice->PaymentTypes")
		}
		result.InvoiceSummary.PaymentTypes = payments
	}

	if req.OtherDocs != nil {
		otherDocs, err := common.MapCommonRequestOtherDocuments(req.OtherDocs)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestOtherDocuments", "MapToInvoiceData", err, "ErrorMapping", "Invoice->OtherDocs")
		}
		result.OtherDocs = otherDocs
	}

	if req.RelatedDocs != nil {
		relatedDocs, err := common.MapCommonRequestRelatedDocuments(req.RelatedDocs)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestRelatedDocuments", "MapToInvoiceData", err, "ErrorMapping", "Invoice->RelatedDocs")
		}
		result.RelatedDocs = relatedDocs
	}

	if req.Appendixes != nil {
		appendixes, err := common.MapCommonRequestAppendix(req.Appendixes)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapAppendixes", "MapToInvoiceData", err, "ErrorMapping", "Invoice->Appendixes")
		}
		result.Appendixes = appendixes
	}

	return nil
}
