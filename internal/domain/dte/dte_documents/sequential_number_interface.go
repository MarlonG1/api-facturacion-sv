package dte_documents

import "context"

// SequentialNumberManager establece los métodos que debe implementar un repositorio de números secuenciales
type SequentialNumberManager interface {
	// GetNextControlNumber obtiene el siguiente número de control
	GetNextControlNumber(ctx context.Context, dteType string, branchID uint, posCode, establishmentCode *string) (string, error)
}
