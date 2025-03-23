package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseReceiver mapea un receptor a un modelo de receptor -> Origen: Response
func MapCommonResponseReceiver(receiver interfaces.Receiver) structs.DTEReceiver {
	result := structs.DTEReceiver{
		TipoDocumento:   receiver.GetDocumentType(),
		NumDocumento:    receiver.GetDocumentNumber(),
		Nombre:          receiver.GetName(),
		NRC:             receiver.GetNRC(),
		NIT:             receiver.GetNIT(),
		Direccion:       addressToPointer(MapCommonResponseAddress(receiver.GetAddress())),
		Correo:          receiver.GetEmail(),
		Telefono:        receiver.GetPhone(),
		CodActividad:    receiver.GetActivityCode(),
		DescActividad:   receiver.GetActivityDescription(),
		NombreComercial: receiver.GetCommercialName(),
	}
	return result
}

func addressToPointer(address structs.DTEAddress) *structs.DTEAddress {
	if address.Departamento == "" && address.Municipio == "" && address.Complemento == "" {
		return nil
	}
	return &address
}
