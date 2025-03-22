package user

import (
	"encoding/json"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/domain/core/error"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"time"
)

type User struct {
	ID                   uint      `json:"-"`
	Status               bool      `json:"-"`
	NIT                  string    `json:"nit"`
	NRC                  string    `json:"nrc"`
	AuthType             string    `json:"auth_type"`
	PasswordPri          string    `json:"password_pri"`
	CommercialName       string    `json:"commercial_name"`
	Business             string    `json:"business_name"`
	EconomicActivity     string    `json:"economic_activity"`
	EconomicActivityDesc string    `json:"economic_activity_desc"`
	Email                string    `json:"email"`
	Phone                string    `json:"phone"`
	YearInDTE            bool      `json:"year_in_dte"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`

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

	if _, err := identification.NewActivityCode(u.EconomicActivity); err != nil {
		return err
	}

	if _, err := base.NewPhone(u.Phone); err != nil {
		return err
	}

	if _, err := base.NewEmail(u.Email); err != nil {
		return err
	}

	if u.AuthType == "" {
		return dte_errors.NewValidationError("RequiredField", "auth_type")
	}

	if u.PasswordPri == "" {
		return dte_errors.NewValidationError("RequiredField", "password_pri")
	}

	if u.CommercialName == "" {
		if len(u.CommercialName) > 150 {
			return dte_errors.NewValidationError("InvalidLength", "commercial_name", "1 to 150", u.CommercialName)
		}
		return dte_errors.NewValidationError("RequiredField", "commercial_name")
	}

	if u.Business == "" {
		if len(u.Business) > 200 {
			return dte_errors.NewValidationError("InvalidLength", "business_name", "1 to 200", u.Business)
		}

		return dte_errors.NewValidationError("RequiredField", "business_name")
	}

	if u.EconomicActivityDesc == "" {
		if len(u.EconomicActivityDesc) > 150 {
			return dte_errors.NewValidationError("InvalidLength", "economic_activity_desc", "1 to 150", u.EconomicActivityDesc)
		}

		return dte_errors.NewValidationError("RequiredField", "economic_activity_desc")
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

	return nil, dte_errors.NewFormattedValidationError(errPackage.ErrBranchMatrixNotFound)
}

// ValidateBranchOffices valida las sucursales del usuario para que cumplan con las reglas de negocio
func (u *User) ValidateBranchOffices() error {
	var matrixCount int
	var matrixHasAddress bool

	// 1. Validar que tenga al menos una sucursal
	if len(u.BranchOffices) == 0 {
		return dte_errors.NewFormattedValidationError(errPackage.ErrAtLeastOneBranch)
	}

	// 2. Validar cada sucursal individualmente
	for _, branchOffice := range u.BranchOffices {
		if err := branchOffice.Validate(); err != nil {
			return err
		}

		if branchOffice.EstablishmentType == constants.CasaMatriz {
			matrixCount++
			if branchOffice.Address != nil {
				matrixHasAddress = true
			}
		}
	}

	// 3. Validar que tenga una casa matriz
	if matrixCount == 0 {
		return dte_errors.NewFormattedValidationError(errPackage.ErrDontHaveBranchMatrix)
	}

	// 4. Validar que la casa matriz tenga dirección
	if !matrixHasAddress {
		return dte_errors.NewFormattedValidationError(errPackage.ErrBranchMatrixWithoutAddress)
	}

	// 5. Validar que tenga solo una casa matriz
	if matrixCount > 1 {
		return dte_errors.NewFormattedValidationError(errPackage.ErrMoreThanOneBranchMatrix)
	}

	return nil
}

// SetBranchesKeysAndSecrets asigna las llaves y secretos a las sucursales del usuario
func (u *User) SetBranchesKeysAndSecrets(keys []string, secrets []string) {
	for i := range u.BranchOffices {
		u.BranchOffices[i].APIKey = keys[i]
		u.BranchOffices[i].APISecret = secrets[i]
		u.BranchOffices[i].IsActive = true
	}
}

func (u *User) ToStringJSON() string {
	jsonUser, _ := json.Marshal(u)
	return string(jsonUser)
}

func (u *User) ListBranches() []ListBranchesResponse {
	var branches []ListBranchesResponse

	for i, branch := range u.BranchOffices {
		branches = append(branches, ListBranchesResponse{
			BranchNumber:      i + 1,
			EstablishmentType: branch.EstablishmentType,
			EstablishmentCode: branch.EstablishmentCode,
			APIKey:            branch.APIKey,
			APISecret:         branch.APISecret,
		})
	}

	return branches
}
