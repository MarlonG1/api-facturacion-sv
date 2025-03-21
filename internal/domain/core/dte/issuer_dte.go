package dte

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
)

type IssuerDTE struct {
	NIT                 string
	NRC                 string
	CommercialName      string
	BusinessName        string
	EstablishmentCode   *string
	EstablishmentCodeMH *string
	Email               *string
	Phone               *string
	EstablishmentType   string
	EstablishmentTypeMH *string
	POSCode             *string
	POSCodeMH           *string
	Address             *user.Address
}
