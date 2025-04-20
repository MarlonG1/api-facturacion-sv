package test

import (
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/sirupsen/logrus"
)

// TestMain configura el ambiente para todas las pruebas de integración
func TestMain(m *testing.T) {
	// Encontrar la raíz del proyecto
	rootPath := utils.FindProjectRoot()

	// Inicializar configuración
	err := config.InitEnvConfig(rootPath)
	if err != nil {
		panic("Error initializing environment config: " + err.Error())
	}

	// Configurar a modo de prueba
	config.Server.AmbientCode = "00"
	config.Server.Debug = true

	// Inicializar el tiempo
	err = utils.TimeInit()
	if err != nil {
		panic("Error initializing time: " + err.Error())
	}

	// Inicializar traducciones
	err = i18n.InitTranslations(rootPath+"/internal/i18n", "en")
	if err != nil {
		panic("Error initializing translations: " + err.Error())
	}

	// Inicializar logger
	err = logs.InitLogger(logrus.DebugLevel.String(), config.Log.Path)
	if err != nil {
		panic("Error initializing logger: " + err.Error())
	}
}
