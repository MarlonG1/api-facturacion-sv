package main

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"os"
)

func main() {
	// Crear e inicializar la aplicación
	app := bootstrap.NewApplication()
	if err := app.Initialize(); err != nil {
		logs.Fatal("Failed to initialize application", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	// Iniciar la aplicación
	if err := app.Start(); err != nil {
		logs.Fatal("Application error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
