package interfaces

import "context"

// ContingencyManager es una interfaz que define los métodos para el manejo de documentos en estado de contingencia.
type ContingencyManager interface {
	// SaveInContingency almacena un documento en la base de datos en estado de contingencia.
	SaveInContingency(ctx context.Context, document interface{}, contingencyType int, reason string) error
}
