package ports

import (
	"context"
)

// DTERepositoryPort es una interfaz que define los métodos de un repositorio de DTE.
type DTERepositoryPort interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error
	Update(ctx context.Context, id, status string, receptionStamp *string) error
}
