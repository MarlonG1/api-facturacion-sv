package user

type ListBranchesResponse struct {
	BranchNumber      int     `json:"branch_number"`
	EstablishmentType string  `json:"establishment_type"`
	EstablishmentCode *string `json:"establishment_code,omitempty"`
	APIKey            string  `json:"api_key"`
	APISecret         string  `json:"api_secret"`
}
