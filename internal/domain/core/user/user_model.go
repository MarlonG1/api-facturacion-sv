package user

import (
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

	return nil
}
