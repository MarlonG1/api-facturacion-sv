package invalidation

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

func MapInvalidationReasonRequest(reason *structs.ReasonRequest) (*models.InvalidationReason, error) {
	invalidationType, err := document.NewInvalidationType(reason.Type)
	if err != nil {
		return nil, err
	}

	responsibleDocType, err := document.NewDTETypeForReceiver(reason.ResponsibleDocType)
	if err != nil {
		return nil, err
	}

	requestorDocType, err := document.NewDTETypeForReceiver(reason.RequestorDocType)
	if err != nil {
		return nil, err
	}

	responsibleName := reason.ResponsibleName
	requestorName := reason.RequestorName

	responsibleDocNum, err := identification.NewDocumentNumber(reason.ResponsibleNumDoc, reason.ResponsibleDocType)
	if err != nil {
		return nil, err
	}

	requestorDocNum, err := identification.NewDocumentNumber(reason.RequestorNumDoc, reason.RequestorDocType)
	if err != nil {
		return nil, err
	}

	result := &models.InvalidationReason{
		Type:               *invalidationType,
		ResponsibleName:    responsibleName,
		ResponsibleDocType: *responsibleDocType,
		ResponsibleDocNum:  *responsibleDocNum,
		RequesterName:      requestorName,
		RequesterDocType:   *requestorDocType,
		RequesterDocNum:    *requestorDocNum,
	}

	if reason.Reason != nil && reason.Type == 3 {
		invalidationReason, err := document.NewInvalidationReason(*reason.Reason)
		if err != nil {
			return nil, err
		}
		result.Reason = invalidationReason
	}

	return result, nil
}
