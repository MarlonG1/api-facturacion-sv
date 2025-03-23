package ports

import "context"

// SequentialNumberRepositoryPort establece los métodos que debe implementar un repositorio de números secuenciales
type SequentialNumberRepositoryPort interface {
	// GetNext obtiene el siguiente número secuencial
	GetNext(ctx context.Context, dteType string, branchID uint) (int, error)
}
