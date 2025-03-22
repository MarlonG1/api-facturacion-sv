package user

import (
	"encoding/json"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/domain/core/error"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"time"
)

type User struct {
	ID             uint      `json:"-"`
	NIT            string    `json:"nit"`
	NRC            string    `json:"nrc"`
	Status         bool      `json:"status"`
	AuthType       string    `json:"auth_type"`
	PasswordPri    string    `json:"password_pri"`
	CommercialName string    `json:"commercial_name"`
	Business       string    `json:"business_name"`
	Email          string    `json:"email"`
	YearInDTE      bool      `json:"year_in_dte"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`

	// Relationships
	BranchOffices []BranchOffice `json:"branch_offices,omitempty"`
}

// Validate válida los campos del usuario para que cumplan con las reglas de negocio
func (u *User) Validate() error {
	if _, err := identification.NewNIT(u.NIT); err != nil {
		return err
	}

	if _, err := identification.NewNRC(u.NRC); err != nil {
		return err
	}

	if u.AuthType == "" {
		return dte_errors.NewValidationError("RequiredField", "auth_type")
	}

	if u.PasswordPri == "" {
		return dte_errors.NewValidationError("RequiredField", "password_pri")
	}

	if u.CommercialName == "" {
		return dte_errors.NewValidationError("RequiredField", "commercial_name")
	}

	if u.Business == "" {
		return dte_errors.NewValidationError("RequiredField", "business_name")
	}

	if u.Email == "" {
		return dte_errors.NewValidationError("RequiredField", "email")
	}

	if u.BranchOffices == nil {
		return dte_errors.NewValidationError("RequiredField", "branch_offices")
	}

	if err := u.ValidateBranchOffices(); err != nil {
		return err
	}

	return nil
}

// GetBranchOfficeMatrix retorna la sucursal que es la casa matriz
func (u *User) GetBranchOfficeMatrix() (*BranchOffice, error) {
	// 1. Buscar la casa matriz
	for _, branchOffice := range u.BranchOffices {
		if branchOffice.EstablishmentType == constants.CasaMatriz {
			return &branchOffice, nil
		}
	}

	return nil, errPackage.ErrBranchMatrixNotFound
}

// ValidateBranchOffices valida las sucursales del usuario para que cumplan con las reglas de negocio
func (u *User) ValidateBranchOffices() error {
	var matrixCount int

	// 1. Validar que tenga al menos una sucursal
	if len(u.BranchOffices) == 0 {
		return errPackage.ErrAtLeastOneBranch
	}

	// 2. Validar cada sucursal individualmente
	for _, branchOffice := range u.BranchOffices {
		if err := branchOffice.Validate(); err != nil {
			return err
		}

		if branchOffice.EstablishmentType == constants.CasaMatriz {
			matrixCount++
		}
	}

	// 3. Validar que tenga una casa matriz
	if matrixCount == 0 {
		return errPackage.ErrDontHaveBranchMatrix
	}

	// 4. Validar que tenga solo una casa matriz
	if matrixCount > 1 {
		return errPackage.ErrMoreThanOneBranchMatrix
	}

	return nil
}

// SetBranchesKeysAndSecrets asigna las llaves y secretos a las sucursales del usuario
func (u *User) SetBranchesKeysAndSecrets(keys []string, secrets []string) {
	for i, branch := range u.BranchOffices {
		branch.APIKey = keys[i]
		branch.APISecret = secrets[i]
		branch.IsActive = true
	}
}

func (u *User) ToStringJSON() string {
	jsonUser, _ := json.Marshal(u)
	return string(jsonUser)
}
