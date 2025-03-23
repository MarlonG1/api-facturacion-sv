package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseIdentification mapea la identificación de una invoice electrónica a un modelo de identificación -> Origen: Response
func MapCommonResponseIdentification(identification interfaces.Identification) *structs.DTEIdentification {
	return &structs.DTEIdentification{
		Version:          identification.GetVersion(),
		Ambiente:         identification.GetAmbient(),
		TipoDte:          identification.GetDTEType(),
		NumeroControl:    identification.GetControlNumber(),
		CodigoGeneracion: identification.GetGenerationCode(),
		TipoModelo:       identification.GetModelType(),
		TipoOperacion:    identification.GetOperationType(),
		FecEmi:           identification.GetEmissionDate().Format("2006-01-02"),
		HorEmi:           identification.GetEmissionTime().Format("15:04:05"),
		TipoMoneda:       identification.GetCurrency(),
	}
}
