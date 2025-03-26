package request_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
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

func (i *InvalidationMapper) MapToInvalidationDocument(req *structs.InvalidationRequest, client *dte.IssuerDTE, baseDte *dte.DTEDetails, emissionDate time.Time) (*models.InvalidationDocument, error) {
	if req == nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationMapper", "MapToInvalidationDocument", "Invalid request", nil)
	}

	reason, err := invalidation.MapInvalidationReasonRequest(req.Reason)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationMapper", "MapToInvalidationDocument", "Error mapping reason", err)
	}

	document, err := invalidation.MapInvalidatedDocument(baseDte, req, emissionDate)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationMapper", "MapToInvalidationDocument", "Error mapping invalidated document", err)
	}

	identification, err := common.MapCommonRequestIdentification(1, 2, document.DocumentType.GetValue())
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationMapper", "MapToInvalidationDocument", "Error mapping identification", err)
	}

	newUUID := uuid.New().String()
	identification.GenerationCode = *identificationVO.NewValidatedGenerationCode(strings.ToUpper(newUUID))

	issuer, err := common.MapCommonIssuer(client)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("CCFMapper", "MapToCCFData", "Error mapping issuer", err)
	}

	return &models.InvalidationDocument{
		Identification: identification,
		Reason:         reason,
		Issuer:         issuer,
		Document:       document,
	}, nil
}

func (i *InvalidationMapper) ValidateInvalidationReRequest(req *structs.InvalidationRequest) error {
	if req == nil {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty request", nil)
	}

	if req.Reason == nil {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty reason", nil)
	}

	if (req.Reason.Type == 1 || req.Reason.Type == 3) && req.ReplacementGenerationCode == nil {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Replacement code is required for type 1 and 3", nil)
	}

	if req.Reason.Type == 2 && req.ReplacementGenerationCode != nil {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Replacement code must be null for type 2", nil)
	}

	if req.GenerationCode == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty generation code", nil)
	}

	if req.Reason.Type == 0 {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty reason type", nil)
	}

	if req.Reason.ResponsibleName == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty responsible name", nil)
	}

	if req.Reason.ResponsibleDocType == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty responsible document type", nil)
	}

	if req.Reason.ResponsibleNumDoc == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty responsible document number", nil)
	}

	if req.Reason.RequestorName == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty requestor name", nil)
	}

	if req.Reason.RequestorDocType == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty requestor document type", nil)
	}

	if req.Reason.RequestorNumDoc == "" {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "Empty requestor document number", nil)
	}

	if req.Reason.Reason == nil && req.Reason.Type == 3 {
		return shared_error.NewGeneralServiceError("InvalidationMapper", "validateInvoiceRequest", "If the reason type is 3, the reason field in reason object must be provided", nil)
	}

	return nil
}
