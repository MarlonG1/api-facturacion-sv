package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"time"
)

// AuthCredentials representa las credenciales de autenticación
type AuthCredentials struct {
	MHCredentials *HaciendaCredentials `json:"credentials"`
	APIKey        string               `json:"api_key"`
	APISecret     string               `json:"api_secret"`
}

func (a *AuthCredentials) Validate() error {
	if a.APIKey == "" {
		return dte_errors.NewValidationError("RequiredField", "api_key")
	}
	if a.APISecret == "" {
		return dte_errors.NewValidationError("RequiredField", "api_secret")
	}

	if a.MHCredentials == nil {
		return dte_errors.NewValidationError("RequiredField", "credentials")
	}

	if a.MHCredentials.Username == "" {
		return dte_errors.NewValidationError("RequiredField", "username")
	}

	if a.MHCredentials.Password == "" {
		return dte_errors.NewValidationError("RequiredField", "password")
	}

	return nil
}

// AuthClaims representa la información que se incluirá en el token JWT
type AuthClaims struct {
	ClientID  uint      `json:"sub"`
	BranchID  uint      `json:"branch_sub"`
	AuthType  string    `json:"auth_type"`
	NIT       string    `json:"nit"`
	ExpiresAt time.Time `json:"expires_at"`
}

// HaciendaCredentials representa las credenciales de hacienda
type HaciendaCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
