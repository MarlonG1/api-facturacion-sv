package models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/location"

// Address es una estructura que representa un Department, Municipality y Complement de un DTE
type Address struct {
	Department   location.Department   `json:"department"`
	Municipality location.Municipality `json:"municipality"`
	Complement   location.Address      `json:"complement"`
}

func (a *Address) GetDepartment() string {
	return a.Department.GetValue()
}

func (a *Address) GetMunicipality() string {
	return a.Municipality.GetValue()
}

func (a *Address) GetComplement() string {
	return a.Complement.GetValue()
}
