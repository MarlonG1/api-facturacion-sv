package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// MapCommonRequestDoctorInfo mapea la información de un doctor a un modelo de doctor -> Origen: Request
func MapCommonRequestDoctorInfo(docInfo structs.DoctorRequest) (*models.DoctorInfo, error) {
	var docNIT *identification.NIT
	var err error
	if docInfo.Name == "" || docInfo.ServiceType == 0 {
		return nil, shared_error.NewGeneralServiceError("CommonMapper", "MapCommonRequestDoctorInfo", "Doctor name and service type are required, but one or both are empty (service_type 0 is considered empty)", nil)
	}

	if docInfo.IdentificationDoc == nil && docInfo.NIT == nil {
		return nil, shared_error.NewGeneralServiceError("CommonMapper", "MapCommonRequestDoctorInfo", "Doctor identification or NIT is required, but both are empty", nil)
	}

	if docInfo.NIT != nil {
		docNIT, err = identification.NewNIT(*docInfo.NIT)
		if err != nil {
			return nil, err
		}
	}

	serviceType, err := document.NewServiceType(docInfo.ServiceType)
	if err != nil {
		return nil, err
	}

	return &models.DoctorInfo{
		Name:           docInfo.Name,
		ServiceType:    *serviceType,
		NIT:            docNIT,
		Identification: docInfo.IdentificationDoc,
	}, nil
}
