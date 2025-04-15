package contingency

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// ContingencyEventSender interfaz para enviar eventos de contingencia
type ContingencyEventSender interface {
	PrepareAndSendContingencyEvent(ctx context.Context, docs []dte.ContingencyDocument) error
}
