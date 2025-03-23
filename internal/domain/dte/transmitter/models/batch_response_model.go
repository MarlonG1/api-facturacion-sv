package models

type BatchResponse struct {
	Version         int     `json:"version"`
	Ambient         string  `json:"ambiente"`
	VersionApp      int     `json:"versionApp"`
	Status          string  `json:"estado"`
	SendID          string  `json:"idEnvio"`
	BatchCode       string  `json:"codigoLote"`
	ProcessingDate  string  `json:"fhProcesamiento"`
	ReceptionStamp  *string `json:"selloRecibido"`
	ClassifyMessage string  `json:"clasificaMsg"`
	MessageCode     string  `json:"codigoMsg"`
	Description     string  `json:"descripcionMsg"`
}

type ConsultBatchResponse struct {
	Processed []HaciendaResponse `json:"procesados"`
	Rejected  []HaciendaResponse `json:"rechazados"`
}
