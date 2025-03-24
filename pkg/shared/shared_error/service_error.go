package shared_error

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
)

type ServiceError struct {
	Type      string
	Operation string
	Message   string
	Err       error
}

func (e *ServiceError) Error() string {

	if e.Err == nil {
		return fmt.Sprintf("%s", e.Message)
	}

	if config.Server.Debug {
		return fmt.Sprintf("[%s] %s: %s | cause: %v", e.Type, e.Operation, e.Message, e.Err)
	}

	return fmt.Sprintf("%s -> %s", e.Message, e.Err)
}

// NewGeneralServiceError crea un nuevo error de servicio general con el tipo de servicio, la operación, el mensaje y el error.
func NewGeneralServiceError(serviceType, op, msg string, err error) *ServiceError {
	return &ServiceError{
		Type:      serviceType,
		Operation: op,
		Message:   msg,
		Err:       err,
	}
}
