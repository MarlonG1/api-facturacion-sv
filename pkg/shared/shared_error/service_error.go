package shared_error

import (
	"errors"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"strings"
)

type ServiceError struct {
	Type      string
	Operation string
	Message   string
	Code      string
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

func NewFormattedGeneralServiceError(serviceType, op, code string, args ...interface{}) *ServiceError {
	code = strings.ToLower(code)
	message := i18n.Translate(fmt.Sprintf("service_errors.%s", code), args)
	return &ServiceError{
		Type:      serviceType,
		Operation: op,
		Message:   message,
		Code:      code,
	}
}

func NewFormattedGeneralServiceWithError(serviceType, op string, err error, code string, args ...interface{}) *ServiceError {
	code = strings.ToLower(code)
	message := i18n.Translate(fmt.Sprintf("service_errors.%s", code), args)
	return &ServiceError{
		Type:      serviceType,
		Operation: op,
		Message:   message,
		Code:      code,
		Err:       err,
	}
}

func (e *ServiceError) GetCode() string {
	if e.Err != nil {
		var dteErr *dte_errors.ValidationError
		if errors.As(e.Err, &dteErr) {
			return dteErr.GetType()
		}
	}

	if e.Code != "" {
		return e.Code
	}

	return strings.ToLower(e.Type)
}
