package interfaces

import "context"

// DTEManager es una interfaz que define los métodos de un administrador de DTE.
type DTEManager interface {
	// Create almacena un DTE en la base de datos con el sello de recepción proporcionado.
	Create(ctx context.Context, document interface{}, receptionStamp *string) error
}
