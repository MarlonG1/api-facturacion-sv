package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
)

type ThirdPartySale struct {
	NIT  identification.NIT `json:"nit"`
	Name string             `json:"name"`
}

func (t ThirdPartySale) GetNIT() string {
	return t.NIT.GetValue()
}

func (t ThirdPartySale) GetName() string {
	return t.Name
}
