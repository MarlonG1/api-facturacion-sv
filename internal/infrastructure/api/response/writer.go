package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/hacienda_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type errorType string

const (
	errorValidation errorType = "VALIDATION" // error de validación de datos (Nivel de creacion en ValueObjects)
	errorBusiness   errorType = "BUSINESS"   // error de negocio (Nivel de creacion en procesos y concordancia de datos)
	errorSystem     errorType = "SYSTEM"     // error de sistema (Nivel de creacion en infraestructura y servicios)
)

type ResponseWriter struct{}

// NewResponseWriter crea una nueva instancia de ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// Success envía una respuesta exitosa con el código de estado y los datos proporcionados.
func (w *ResponseWriter) Success(rw http.ResponseWriter, status int, data interface{}, options *SuccessOptions) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	apiResponse := &APIResponse{
		Success: true,
		Data:    data,
	}

	if options == nil {
		json.NewEncoder(rw).Encode(apiResponse)
		return
	}

	qrLink := GenerateQRLink(options.Ambient, options.GenerationCode, options.EmissionDate)

	json.NewEncoder(rw).Encode(APIDTEResponse{
		Success:        true,
		ReceptionStamp: options.ReceptionStamp,
		QRLink:         &qrLink,
		Data:           data,
	})
}

// Error envía una respuesta de error con el código de estado y el mensaje proporcionado.
func (w *ResponseWriter) Error(rw http.ResponseWriter, status int, message string, details []string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(APIErrorResponse{
		Error: &APIError{
			Message: message,
			Details: details,
			Code:    deriveErrorCode(status),
		},
	})
}

// HandleError maneja los diferentes tipos de errores y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) HandleError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json")

	logs.Error("Error processing request", map[string]interface{}{
		"error_type": getErrorType(err),
		"error":      err.Error(),
	})

	switch errorType := getErrorType(err); errorType {
	case errorValidation:
		w.handleValidationError(rw, err)
	case errorBusiness:
		w.handleBusinessError(rw, err)
	default:
		w.handleSystemError(rw, err)
	}
}

// GenerateQRLink Genera un link para consultar la invoice en la página de la Hacienda
func GenerateQRLink(ambiente, codGeneracion string, fechaEmision time.Time) string {
	return fmt.Sprintf("https://admin.factura.gob.sv/consultaPublica?ambiente=%s&codGen=%s&fechaEmi=%s",
		ambiente, codGeneracion, fechaEmision.Format("2006-01-02"))
}

// handleValidationError maneja los errores de validación y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleValidationError(rw http.ResponseWriter, err error) {
	var dteErr *dte_errors.DTEError
	if errors.As(err, &dteErr) {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: dteErr.Message,
				Details: dteErr.GetValidationErrorsString(),
				Code:    "VALIDATION_ERROR",
			},
		})
		return
	}

	var validationErr *dte_errors.ValidationError
	if errors.As(err, &validationErr) {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: validationErr.Error(),
				Details: []string{i18n.Translate("service_errors.NoDetailsAvailable")},
				Code:    strings.ToUpper(validationErr.GetType()),
			},
		})
	}
}

// handleBusinessError maneja los errores de negocio y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleBusinessError(rw http.ResponseWriter, err error) {
	var haciendaErr *hacienda_error.HaciendaResponseError
	if errors.As(err, &haciendaErr) {
		details := []string{
			fmt.Sprintf("State: %s", haciendaErr.Status),
			fmt.Sprintf("Processed at: %s", haciendaErr.ProcessedAt),
		}

		if len(haciendaErr.Observations) > 0 {
			details = append(details, haciendaErr.Observations...)
		}

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: haciendaErr.Description,
				Code:    fmt.Sprintf("HACIENDA_%s", haciendaErr.Code),
				Details: details,
			},
		})
		return
	}

	var svcErr *shared_error.ServiceError
	if errors.As(err, &svcErr) {
		var details []string
		rw.WriteHeader(http.StatusBadRequest)

		if svcErr.Err != nil {
			details = append(details, svcErr.Err.Error())
		} else {
			details = append(details, i18n.Translate("service_errors.NoDetailsAvailable"))
		}

		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: svcErr.Message,
				Details: details,
				Code:    strings.ToUpper(svcErr.GetCode()),
			},
		})
		return
	}

	var httpErr *hacienda_error.HTTPResponseError
	if errors.As(err, &httpErr) {
		rw.WriteHeader(httpErr.StatusCode)
		detail := i18n.Translate("service_errors.ContingencyActiveTransmission")
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: httpErr.Error(),
				Details: []string{detail},
				Code:    fmt.Sprintf("HACIENDA_%d", httpErr.StatusCode),
			},
		})
		return
	}

	rw.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(rw).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Message: err.Error(),
			Code:    "BUSINESS_ERROR",
		},
	})
}

// handleSystemError maneja los errores de sistema y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleSystemError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	message := i18n.Translate("validation_errors.ServerError")
	json.NewEncoder(rw).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Message: message,
			Code:    "SYSTEM_ERROR",
		},
	})
}

// deriveErrorCode deriva el código de error de acuerdo al estado y mensaje proporcionado.
func deriveErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

// getErrorType obtiene el tipo de error de acuerdo al tipo de error proporcionado.
func getErrorType(err error) errorType {
	switch err.(type) {
	case *dte_errors.DTEError, *dte_errors.ValidationError:
		return errorValidation
	case *shared_error.ServiceError, *hacienda_error.HaciendaResponseError, *hacienda_error.HTTPResponseError:
		return errorBusiness
	default:
		return errorSystem
	}
}
