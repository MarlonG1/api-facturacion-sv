package transmitter

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/models"
)

type DocumentProcessor interface {
	ProcessRequest(signedDoc string, document interface{}) (*models.HaciendaRequest, error)
	ProcessResponse(resp *models.HaciendaResponse) (*models.TransmitResult, error)
}
