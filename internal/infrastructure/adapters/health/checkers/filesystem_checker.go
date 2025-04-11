package checkers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/dimiro1/health"
	"os"
	"path/filepath"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type fileSystemChecker struct {
	logPath string
}

func NewFileSystemChecker() ports.ComponentChecker {
	// Obtener la ruta absoluta del proyecto
	logFilePath := filepath.Join(utils.FindProjectRoot()+config.Log.Path, "dte_microservice.log")

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
	// CustomHealthChecker para el sistema de archivos
	health := c.checkHealth()

	status := constants.StatusUp
	details := "File system is healthy"

	if health.IsDown() {
		status = constants.StatusDown
		details = "File system is down"

		// Extraer detalles si están disponibles
		if health.GetInfo("error") != nil {
			details = fmt.Sprintf("%s: %v", details, health.GetInfo("error"))
		}
	}

	return models.Health{
		Status:  status,
		Details: details,
	}
}

func (c *fileSystemChecker) checkHealth() health.Health {
	result := health.NewHealth()

	if err := c.checkFileSystem(); err != nil {
		result.Down()
		result.AddInfo("error", err.Error())
	}

	result.Up()
	return result
}

func (c *fileSystemChecker) checkFileSystem() error {
	// Aseguramos que el directorio existe
	dir := filepath.Dir(c.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory")
	}

	// Verificamos permisos de escritura
	file, err := os.OpenFile(c.logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("the system don't have permissions to write in the log file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logs.Logger.Error("Error closing file:", err)
		}
	}(file)

	return nil
}
