package response

import "time"

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIListResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type APIErrorResponse struct {
	Error *APIError `json:"error"`
}

type APIDTEResponse struct {
	Success        bool        `json:"success"`
	ReceptionStamp *string     `json:"reception_stamp"`
	QRLink         *string     `json:"qr_link"`
	Data           interface{} `json:"data"`
}

type SuccessOptions struct {
	Ambient        string
	GenerationCode string
	EmissionDate   time.Time
	ReceptionStamp *string
}

type APIError struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
	Code    string   `json:"code"`
}

type SuccessEndpoint struct {
	Message  string `json:"message"`
	ClientID string `json:"client_id"`
	NIT      string `json:"nit"`
}
