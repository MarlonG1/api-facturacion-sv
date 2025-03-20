package models

type HaciendaResponse struct {
	Version            int      `json:"version"`
	Ambient            string   `json:"ambiente"`
	VersionApp         int      `json:"versionApp"`
	Status             string   `json:"estado"`
	GenerationCode     string   `json:"codigoGeneracion"`
	ReceptionStamp     string   `json:"selloRecibido"`
	ProcessingDate     string   `json:"fhProcesamiento"`
	ClassifyMessage    string   `json:"clasificaMsg"`
	MessageCode        string   `json:"codigoMsg"`
	DescriptionMessage string   `json:"descripcionMsg"`
	Observations       []string `json:"observaciones,omitempty"`
}
