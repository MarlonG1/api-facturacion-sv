package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/ccf"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type CCFMapper struct{}

func NewCCFMapper() *CCFMapper {
	return &CCFMapper{}
}

// MapToCCFData convierte una solicitud de Comprobante de Crédito Fiscal a datos de CCF.
func (m *CCFMapper) MapToCCFData(req *structs.CreateCreditFiscalRequest, client *dte.IssuerDTE) (*ccf_models.CCFData, error) {
	if err := validateCCFRequest(req); err != nil {
		return nil, err
	}

	items, err := ccf.MapCCFItems(req.Items)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping items", err)
	}

	receiver, err := ccf.MapCCFRequestReceiver(req.Receiver)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping receiver", err)
	}

	if req.Receiver.CommercialName == nil {
		return nil, dte_errors.NewValidationError("RequiredField", "Request->Receiver->CommercialName")
	}
	receiver.CommercialName = req.Receiver.CommercialName

	identification, err := common.MapCommonRequestIdentification(constants.ModeloFacturacionPrevio, 3, constants.CCFElectronico)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping identification", err)
	}

	summary, err := ccf.MapCCFRequestSummary(req.Summary)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping summary", err)
	}

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping issuer", err)
	}

	result := &ccf_models.CCFData{
		InputDataCommon: &models.InputDataCommon{
			Issuer:         issuer,
			Identification: identification,
			Receiver:       receiver,
		},
		Items:         items,
		CreditSummary: summary,
	}

	if err = mapCCFOptionalFields(req, result); err != nil {
		return nil, err
	}

	return result, nil
}

// validateCCFRequest valida que la solicitud de Comprobante de Crédito Fiscal sea correcta.
func validateCCFRequest(req *structs.CreateCreditFiscalRequest) error {
	if req == nil {
		return shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "The request is empty", nil)
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
	if req.Receiver.CommercialName == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Receiver->CommercialName")
	}

	if req.Summary.TotalIVA != 0 {
		return dte_errors.NewValidationError("InvalidField", "Request->Summary->TotalIVA")
	}
	if req.Receiver.DocumentType != nil {
		return dte_errors.NewValidationError("InvalidField", "Request->Receiver->DocumentType")
	}
	if req.Receiver.DocumentNumber != nil {
		return dte_errors.NewValidationError("InvalidField", "Request->Receiver->DocumentNumber")
	}
	return nil
}

// mapCCFOptionalFields mapea los campos opcionales de la solicitud de Comprobante de Crédito Fiscal.
func mapCCFOptionalFields(req *structs.CreateCreditFiscalRequest, result *ccf_models.CCFData) error {
	if req.ThirdPartySale != nil {
		thirdPartySale, err := common.MapCommonRequestThirdPartySale(req.ThirdPartySale)
		if err != nil {
			return shared_error.NewGeneralServiceError("MapCommonRequestThirdPartySale", "MapToCCFData", "Error mapping third party sale", err)
		}
		result.ThirdPartySale = thirdPartySale
	}

	if req.Extension != nil {
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
		result.CreditSummary.PaymentTypes = payments
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
