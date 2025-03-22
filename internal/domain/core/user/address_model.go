package user

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
)

// Address representa la dirección de una sucursal o casa matriz
type Address struct {
	ID           uint   `json:"-"`
	BranchID     uint   `json:"-"`
	Municipality string `json:"municipality"`
	Department   string `json:"department"`
	Complement   string `json:"complement"`
}

func (a *Address) Validate() error {
	if a.Municipality == "" {
		return dte_errors.NewValidationError("RequiredField", "municipality")
	}
	if a.Department == "" {
		return dte_errors.NewValidationError("RequiredField", "department")
	}

	if a.Complement == "" {
		return dte_errors.NewValidationError("RequiredField", "complement")
	}

	return nil
}
