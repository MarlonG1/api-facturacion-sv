package models

// DTEDetails representa los detalles de un documento tributario electrónico
type DTEDetails struct {
	ID             string `json:"id,omitempty"`
	DTEType        string `json:"dte_type"`
	ControlNumber  string `json:"control_number"`
	ReceptionStamp string `json:"reception_stamp"`
	Status         string `json:"status"`
	JSONData       string `json:"json_data"`
}
