package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseThirdPartySale mapea una venta de terceros a un modelo de venta de terceros -> Origen: Response
func MapCommonResponseThirdPartySale(sale interfaces.ThirdPartySale) *structs.DTEThirdPartySale {
	if sale == nil {
		return nil
	}

	return &structs.DTEThirdPartySale{
		NIT:    sale.GetNIT(),
		Nombre: sale.GetName(),
	}
}
