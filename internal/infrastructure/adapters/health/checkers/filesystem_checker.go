package checkers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type fileSystemChecker struct {
	logPath string
}

func NewFileSystemChecker() ports.ComponentChecker {
	// Obtener la ruta absoluta del proyecto
	logFilePath := filepath.Join(utils.FindProjectRoot()+env.Log.Path, "dte_microservice.log")

	return &fileSystemChecker{
		logPath: logFilePath,
	}
}

func (c *fileSystemChecker) Name() string {
	return "filesystem"
}

// Check verifica si el sistema de archivos tiene permisos de escritura
// intentando escribir en el archivo de logs
// Devuelve un estado de salud con el resultado de la verificación
func (c *fileSystemChecker) Check() models.Health {
	// Aseguramos que el directorio existe
	dir := filepath.Dir(c.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logs.Error("Failed to create log directory", map[string]interface{}{
			"error": err.Error(),
			"path":  dir,
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: fmt.Sprintf("Failed to create log directory"),
		}
	}

	// Verificamos permisos de escritura
	file, err := os.OpenFile(c.logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logs.Error("Filesystem check failed", map[string]interface{}{
			"error": err.Error(),
			"path":  c.logPath,
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: fmt.Sprintf("Storage system unavailable"),
		}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logs.Error("Failed to close file", map[string]interface{}{
				"error": err.Error(),
				"path":  c.logPath,
			})
		}
	}(file)

	return models.Health{
		Status:  constants.StatusUp,
		Details: fmt.Sprintf("Storage system available"),
	}
}
