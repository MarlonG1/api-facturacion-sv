package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseExtension mapea una extensión a un modelo de extensión -> Origen: Response
func MapCommonResponseExtension(extension interfaces.Extension) *structs.DTEExtension {
	if extension == nil {
		return nil
	}

	return &structs.DTEExtension{
		NombreEntrega:    extension.GetDeliveryName(),
		DocumentoEntrega: extension.GetDeliveryDocument(),
		NombreRecibe:     extension.GetReceiverName(),
		DocumentoRecibe:  extension.GetReceiverDocument(),
		Observacion:      extension.GetObservation(),
		PlacaVehiculo:    extension.GetVehiculePlate(),
	}
}
