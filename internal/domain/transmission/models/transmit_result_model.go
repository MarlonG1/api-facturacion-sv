package models

type TransmitResult struct {
	Status         string
	ReceptionStamp *string
	ProcessingDate string
	MessageCode    string
	MessageDesc    string
	Observations   []string
}
