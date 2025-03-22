package ports

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/models"
)

type BaseTransmitter interface {
	RetryTransmission(context.Context, interface{}, string, string) (*models.TransmitResult, error)
	CheckStatus(context.Context, interface{}, string) (*models.TransmitResult, error)
}
