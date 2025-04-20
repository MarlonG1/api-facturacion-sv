package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// MapCommonRequestThirdPartySale mapea una venta a tercero a un modelo
func MapCommonRequestThirdPartySale(sale *structs.ThirdPartySaleRequest) (*models.ThirdPartySale, error) {
	if sale.Name == "" || sale.NIT == "" {
		return nil, shared_error.NewFormattedGeneralServiceError(
			"CommonMapper",
			"MapCommonRequestThirdPartySale",
			"InvalidThirdParty",
		)
	}

	nit, err := identification.NewNIT(sale.NIT)
	if err != nil {
		return nil, err
	}

	return &models.ThirdPartySale{
		Name: sale.Name,
		NIT:  *nit,
	}, nil
}
