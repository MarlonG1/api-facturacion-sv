package invalidation

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

func MapInvalidationReasonResponse(reason *models.InvalidationReason) *structs.ReasonResponse {
	if reason == nil {
		return nil
	}

	result := &structs.ReasonResponse{
		TipoAnulacion:     reason.Type.GetValue(),
		NombreResponsable: reason.ResponsibleName,
		TipDocResponsable: reason.ResponsibleDocType.GetValue(),
		NumDocResponsable: reason.ResponsibleDocNum.GetValue(),
		NombreSolicita:    reason.RequesterName,
		TipDocSolicita:    reason.RequesterDocType.GetValue(),
		NumDocSolicita:    reason.RequesterDocNum.GetValue(),
	}

	if reason.Reason != nil {
		result.MotivoAnulacion = utils.ToStringPointer(reason.Reason.GetValue())
	}

	return result
}
