package ports

import (
	"context"
	models2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
)

// DTETransmitter define el comportamiento de un transmisor de documentos electrónicos
type DTETransmitter interface {
	// Transmit envía un documento tributario electrónico a Hacienda
	Transmit(context.Context, interface{}, string, string) (*models2.TransmitResult, error)
	// CheckDocumentStatus verifica el estado de un documento tributario electrónico en Hacienda
	CheckDocumentStatus(context.Context, interface{}, string) (*models2.TransmitResult, error)
	// SendToHacienda envía un documento a Hacienda
	SendToHacienda(context.Context, *models2.HaciendaRequest, string) (*models2.HaciendaResponse, error)
}
