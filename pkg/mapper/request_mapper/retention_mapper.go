package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/retention"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type RetentionMapper struct{}

func NewRetentionMapper() *RetentionMapper {
	return &RetentionMapper{}
}

// MapToRetentionData convierte una solicitud de retención a datos de retención_models.
func (m *RetentionMapper) MapToRetentionData(req *structs.CreateRetentionRequest, client *dte.IssuerDTE) (*retention_models.InputRetentionData, error) {
	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("RetentionMapper", "MapToRetentionData", err, "ErrorMapping", "Retention->Issuer")
	}

	items, err := retention.MapRetentionItemList(req.Items)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("RetentionMapper", "MapToRetentionData", err, "ErrorMapping", "Retention->Items")
	}

	identification, err := common.MapCommonRequestIdentification(constants.ModeloFacturacionPrevio, 1, constants.ComprobanteRetencionElectronico)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("RetentionMapper", "MapToRetentionData", err, "ErrorMapping", "Retention->Identification")
	}

	if err = validateReceiverRequest(req); err != nil {
		return nil, err
	}

	receiver, err := retention.MapRetentionRequestReceiver(req.Receiver)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("RetentionMapper", "MapToRetentionData", err, "ErrorMapping", "Retention->Receiver")
	}

	result := &retention_models.InputRetentionData{
		InputDataCommon: &models.InputDataCommon{
			Issuer:         issuer,
			Receiver:       receiver,
			Identification: identification,
		},
		RetentionItems:   items,
		RetentionSummary: &retention_models.RetentionSummary{},
	}

	if err = mapRetentionOptionalFields(req, result); err != nil {
		return nil, err
	}

	if m.IsAllPhysical(req.Items) {

		if req.Summary == nil {
			return nil, dte_errors.NewValidationError("RequiredField", "Request->Summary")
		}
		result.RetentionSummary, err = retention.MapRetentionSummary(req.Summary)
		if err != nil {
			return nil, shared_error.NewFormattedGeneralServiceWithError("RetentionMapper", "MapToRetentionData", err, "ErrorMapping", "Retention->Summary")
		}
	}

	return result, nil
}

func validateReceiverRequest(req *structs.CreateRetentionRequest) error {

	if req.Receiver == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Receiver")
	}

	if req.Receiver.DocumentType == nil {
		return dte_errors.NewValidationError("InvalidField", "Request->Receiver->DocumentType")
	}
	if req.Receiver.DocumentNumber == nil {
		return dte_errors.NewValidationError("InvalidField", "Request->Receiver->DocumentNumber")
	}

	return nil
}

func mapRetentionOptionalFields(req *structs.CreateRetentionRequest, result *retention_models.InputRetentionData) error {
	if req.Extension != nil {
		if req.Extension.VehiculePlate != nil {
			return dte_errors.NewValidationError("InvalidFieldValue", "Request->Extension->VehiculePlate")
		}

		extension, err := common.MapCommonRequestExtension(req.Extension)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapCommonRequestExtension", "MapToRetentionData", err, "ErrorMapping", "Retention->Extension")
		}
		result.Extension = extension
	}

	if req.Appendixes != nil {
		appendixes, err := common.MapCommonRequestAppendix(req.Appendixes)
		if err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("MapAppendixes", "MapToRetentionData", err, "ErrorMapping", "Retention->Appendixes")
		}
		result.Appendixes = appendixes
	}

	return nil
}

func (m *RetentionMapper) IsAllPhysical(items []structs.RetentionItem) bool {
	isAllPhysical := true
	for _, item := range items {
		if item.DocumentType == 2 {
			isAllPhysical = false
		}
	}
	return isAllPhysical
}
