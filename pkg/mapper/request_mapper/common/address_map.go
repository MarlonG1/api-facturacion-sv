package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/location"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// MapCommonRequestAddress mapea una dirección común a un modelo de dirección -> Origen: Request
func MapCommonRequestAddress(address structs.AddressRequest) (*models.Address, error) {
	if address.Department == "" || address.Municipality == "" || address.Complement == "" {
		return nil, shared_error.NewGeneralServiceError("CommonMapper", "MapCommonRequestAddress", "When address is present, the fields department, municipality and complement must be present", nil)
	}

	department, err := location.NewDepartment(address.Department)
	if err != nil {
		return nil, err
	}

	municipality, err := location.NewMunicipality(address.Municipality, *department)
	if err != nil {
		return nil, err
	}

	complement, err := location.NewAddress(address.Complement)
	if err != nil {
		return nil, err
	}

	return &models.Address{
		Department:   *department,
		Municipality: *municipality,
		Complement:   *complement,
	}, nil
}

// MapClientAddress mapea una dirección de cliente a un modelo de dirección -> Origen: Base de datos
func MapClientAddress(address *user.Address) (*models.Address, error) {
	return &models.Address{
		Department:   *location.NewValidatedDepartment(address.Department),
		Municipality: *location.NewValidatedMunicipality(address.Municipality, address.Department),
		Complement:   *location.NewValidatedAddress(address.Complement),
	}, nil
}
