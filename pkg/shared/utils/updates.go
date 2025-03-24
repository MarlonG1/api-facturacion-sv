package utils

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// UpdateContingencyIdentification actualiza la identificación de contingencia en el JSON del DTE.
func UpdateContingencyIdentification(identification *structs.DTEIdentification, contiType *int8, reason *string) {
	identification.TipoModelo = constants.ModeloFacturacionDiferido
	identification.TipoOperacion = constants.TransmisionContingencia
	tipoContingencia := int(*contiType)
	identification.TipoContingencia = &tipoContingencia
	identification.MotivoContin = reason
}
