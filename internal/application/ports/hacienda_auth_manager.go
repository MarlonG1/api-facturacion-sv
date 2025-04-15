package ports

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
)

type HaciendaAuthManager interface {
	// GetOrCreateHaciendaToken obtiene un token de Hacienda, primero verificando la caché
	GetOrCreateHaciendaToken(ctx context.Context, systemToken string) (string, error)
	// GetOrCreateHaciendaTokenWithCreds obtiene un token de Hacienda, primero verificando la caché
	GetOrCreateHaciendaTokenWithCreds(ctx context.Context, systemToken string, creds models.HaciendaCredentials) (string, error)
}
