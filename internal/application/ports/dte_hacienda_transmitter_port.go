package ports

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/models"
)

// DTETransmitter define el comportamiento de un transmisor de documentos electrónicos
type DTETransmitter interface {
	// Transmit envía un documento tributario electrónico a Hacienda
	Transmit(context.Context, interface{}, string, string) (*models.TransmitResult, error)
	// CheckDocumentStatus verifica el estado de un documento tributario electrónico en Hacienda
	CheckDocumentStatus(context.Context, interface{}, string) (*models.TransmitResult, error)
	// SendToHacienda envía un documento a Hacienda
	SendToHacienda(context.Context, *models.HaciendaRequest, string) (*models.HaciendaResponse, error)
}
