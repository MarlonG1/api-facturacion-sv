package models

type BatchRequest struct {
	Version   int      `json:"version"`
	Ambient   string   `json:"ambiente"`
	SendID    string   `json:"idEnvio"`
	NIT       string   `json:"nitEmisor"`
	Documents []string `json:"documentos"`
}
