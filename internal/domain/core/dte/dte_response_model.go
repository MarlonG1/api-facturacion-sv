package dte

type DTEResponse struct {
	ControlNumber  string                 `json:"control_number"`
	GenerationCode string                 `json:"generation_code"`
	ReceptionStamp *string                `json:"reception_stamp"`
	Status         string                 `json:"status"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
	JSONData       map[string]interface{} `json:"json_data"`
}
