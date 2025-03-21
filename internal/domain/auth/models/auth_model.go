package models

import "time"

// AuthCredentials representa las credenciales de autenticación
type AuthCredentials struct {
	AuthType      string                 `json:"auth_type"`
	MHCredentials *HaciendaCredentials   `json:"credentials"`
	APIKey        string                 `json:"api_key,omitempty"`
	APISecret     string                 `json:"api_secret,omitempty"`
	VaultConfig   map[string]interface{} `json:"vault_config,omitempty"`
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
