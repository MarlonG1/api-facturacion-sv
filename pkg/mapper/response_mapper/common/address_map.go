package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseAddress mapea una dirección a una dirección DTE
func MapCommonResponseAddress(address interfaces.Address) structs.DTEAddress {
	return structs.DTEAddress{
		Departamento: address.GetDepartment(),
		Municipio:    address.GetMunicipality(),
		Complemento:  address.GetComplement(),
	}
}
