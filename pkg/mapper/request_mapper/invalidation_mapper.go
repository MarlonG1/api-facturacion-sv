package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	identificationVO "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/invalidation"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/google/uuid"
	"strings"
	"time"
)

type InvalidationMapper struct{}

func NewInvalidationMapper() *InvalidationMapper {
	return &InvalidationMapper{}
}

func (i *InvalidationMapper) MapToInvalidationData(req *structs.CreateInvalidationRequest, client *dte.IssuerDTE, baseDte *dte.DTEDetails, emissionDate time.Time) (*models.InvalidationDocument, error) {
	if req == nil {
		return nil, dte_errors.NewValidationError("RequiredField", "Request")
	}

	reason, err := invalidation.MapInvalidationReasonRequest(req.Reason)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvalidationMapper", "MapToInvalidationData", err, "ErrorMapping", "Invalidation->Reason")
	}

	document, err := invalidation.MapInvalidatedDocument(baseDte, req, emissionDate)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvalidationMapper", "MapToInvalidationData", err, "ErrorMapping", "Invalidation->Document")
	}

	identification, err := common.MapCommonRequestIdentification(1, 2, document.DocumentType.GetValue())
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvalidationMapper", "MapToInvalidationData", err, "ErrorMapping", "Invalidation->Identification")
	}

	newUUID := uuid.New().String()
	identification.GenerationCode = *identificationVO.NewValidatedGenerationCode(strings.ToUpper(newUUID))

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("InvalidationMapper", "MapToInvalidationData", err, "ErrorMapping", "Invalidation->Issuer")
	}

	return &models.InvalidationDocument{
		Identification: identification,
		Reason:         reason,
		Issuer:         issuer,
		Document:       document,
	}, nil
}

func (i *InvalidationMapper) ValidateInvalidationReRequest(req *structs.CreateInvalidationRequest) error {
	if req == nil {
		return dte_errors.NewValidationError("RequiredField", "Request")
	}

	if req.Reason == nil {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason")
	}

	if (req.Reason.Type == 1 || req.Reason.Type == 3) && req.ReplacementGenerationCode == nil {
		return shared_error.NewFormattedGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "InvalidReplacementCode1And3")
	}

	if req.Reason.Type == 2 && req.ReplacementGenerationCode != nil {
		return shared_error.NewFormattedGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "InvalidInvalidationType2")
	}

	if req.GenerationCode == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->GenerationCode")
	}

	if req.Reason.Type == 0 {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->Type")
	}

	if req.Reason.ResponsibleName == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->ResponsibleName")
	}

	if req.Reason.ResponsibleDocType == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->ResponsibleDocType")
	}

	if req.Reason.ResponsibleNumDoc == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->ResponsibleNumDoc")
	}

	if req.Reason.RequestorName == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->RequestorName")
	}

	if req.Reason.RequestorDocType == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->RequestorDocType")
	}

	if req.Reason.RequestorNumDoc == "" {
		return dte_errors.NewValidationError("RequiredField", "Request->Reason->RequestorNumDoc")
	}

	if req.Reason.Reason == nil && req.Reason.Type == 3 {
		return shared_error.NewFormattedGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "InvalidInvalidationType3")
	}

	return nil
}
