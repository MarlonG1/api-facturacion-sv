package ports

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"

// CircuitManager define una interfaz para implementaciones de circuit breaker
type CircuitManager interface {
	// AllowRequest determina si una solicitud debe ser permitida basada en el estado actual
	AllowRequest() bool
	// RecordSuccess registra una operación exitosa
	RecordSuccess()
	// RecordFailure registra un fallo en la operación
	RecordFailure()
	// GetState devuelve el estado actual del circuit breaker
	GetState() constants.State
	// GetFailureCount devuelve el número actual de fallos registrados
	GetFailureCount() int32
}
