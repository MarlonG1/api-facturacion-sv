package models

type BranchOffice struct {
	ID                  uint    `json:"id,omitempty"`
	UserID              uint    `json:"user_id"`
	EstablishmentCode   *string `json:"establishment_code,omitempty"`
	Email               *string `json:"email,omitempty"`
	APIKey              string  `json:"api_key"`
	APISecret           string  `json:"api_secret"`
	Phone               *string `json:"phone,omitempty"`
	EstablishmentType   string  `json:"establishment_type"`
	EstablishmentTypeMH *string `json:"establishment_type_mh,omitempty"`
	POSCode             *string `json:"pos_code,omitempty"`
	POSCodeMH           *string `json:"pos_code_mh,omitempty"`
	IsActive            bool    `json:"is_active"`
}
