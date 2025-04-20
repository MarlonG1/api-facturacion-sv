package retention

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapRetentionResponseExtension mapea una extensión a un modelo de extensión -> Origen: Response
func MapRetentionResponseExtension(extension interfaces.Extension) *structs.RetentionExtension {
	if extension == nil {
		return nil
	}

	return &structs.RetentionExtension{
		NombreEntrega:    extension.GetDeliveryName(),
		DocumentoEntrega: extension.GetDeliveryDocument(),
		NombreRecibe:     extension.GetReceiverName(),
		DocumentoRecibe:  extension.GetReceiverDocument(),
		Observacion:      extension.GetObservation(),
	}
}
