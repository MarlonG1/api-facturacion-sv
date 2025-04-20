package constants

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
)

// GetErrorMessage obtiene el mensaje de error según el idioma configurado
func GetErrorMessage(errorCode string, params ...interface{}) string {
	message := i18n.Translate(fmt.Sprintf("validation_errors.%s", errorCode), params...)

	// Si el modo debug está activado, se muestra el código de error
	if config.Server.Debug {
		return fmt.Sprintf("[%s] %s", errorCode, message)
	}

	return message
}
