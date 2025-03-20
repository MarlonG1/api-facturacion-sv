package models

type HaciendaRequest struct {
	Ambient        string      `json:"ambiente"`
	SendID         int         `json:"idEnvio"`
	Version        int         `json:"version"`
	Document       interface{} `json:"documento"`
	DTEType        string      `json:"tipoDte"`
	GenerationCode string      `json:"codigoGeneracion"`
	URL            string      `json:"-"`
}
