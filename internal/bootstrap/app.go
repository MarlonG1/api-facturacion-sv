// bootstrap/app.go
package bootstrap

import (
	"context"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap/containers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/cmd/setup"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	errPackage "github.com/MarlonG1/api-facturacion-sv/config/error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/server"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// Application representa la aplicación completa
type Application struct {
	server       *server.Server
	container    *containers.Container
	dbConnection *drivers.DbConnection
}

// SupportedDrivers contiene la configuración de drivers de base de datos soportados
var SupportedDrivers = map[string]drivers.DriverConfig{
	"mysql":    drivers.NewMysqlDriver(),
	"postgres": drivers.NewPostgresDriver(),
}

// NewApplication crea una nueva instancia de la aplicación
func NewApplication() *Application {
	return &Application{}
}

// Initialize inicializa todos los componentes de la aplicación
func (app *Application) Initialize() error {
	// 0. Obtener el root path del proyecto
	rootPath := utils.FindProjectRoot()

	// 1. Inicializar la configuración del entorno
	err := config.InitEnvConfig(rootPath)
	if err != nil {
		return fmt.Errorf("error initializing environment configuration: %w", err)
	}

	// 2. Inicializar el logger
	err = logs.InitLogger(config.Log.Level, config.Log.Path)
	if err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	logs.Info("Logger initialized successfully")

	// 3. Inicializar el tiempo global
	err = utils.TimeInit()
	if err != nil {
		logs.Fatal("Failed to initialize global time", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("error initializing global time: %w", err)
	}

	// 4. Iniciar la configuración de la base de datos y las migraciones
	app.dbConnection, err = app.initDatabaseConfigurations()
	if err != nil {
		logs.Fatal("Failed to initialize database configurations", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("error initializing database configurations: %w", err)
	}

	// 5. Inicializar el contenedor de dependencias
	app.container = containers.NewContainer(app.dbConnection)
	err = app.container.Initialize()
	if err != nil {
		logs.Error("Failed to initialize container", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("error initializing container: %w", err)
	}

	// 6. Inicializar el servidor
	app.server = server.Initialize(app.container)

	// 7. Inicializar los jobs
	err = setup.SetupJobs(app.container.Services().ContingencyManager(), config.Server.AmbientCode, app.dbConnection)
	if err != nil {
		logs.Error("Failed to setup jobs", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("error setting up jobs: %w", err)
	}

	return nil
}

// Start inicia la aplicación y maneja señales para un apagado controlado
func (app *Application) Start() error {
	// Canal para recibir señales del sistema operativo
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar el servidor en una goroutine
	go func() {
		logs.Info("Server started successfully", map[string]interface{}{"port": config.Server.Port})
		serverErrors <- app.server.Start()
	}()

	// Esperar por señales o errores
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-signals:
		logs.Info("Shutdown signal received", map[string]interface{}{"signal": sig.String()})

		// Crear un contexto con timeout para el apagado controlado
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Cerrar el servidor HTTP de forma controlada
		if err := app.server.Shutdown(ctx); err != nil {
			logs.Error("Server shutdown error", map[string]interface{}{"error": err.Error()})
			return fmt.Errorf("server shutdown error: %w", err)
		}

		// Cerrar la conexión a la base de datos
		if err := app.dbConnection.Close(); err != nil {
			logs.Error("Database connection close error", map[string]interface{}{"error": err.Error()})
			return fmt.Errorf("database connection close error: %w", err)
		}

		logs.Info("Shutdown completed", nil)
	}

	return nil
}

// initDatabaseConfigurations inicializa las configuraciones de la base de datos
func (app *Application) initDatabaseConfigurations() (*drivers.DbConnection, error) {
	// 1. Seleccionar el driver de la base de datos
	driver := app.selectDatabaseDriver()
	if driver == nil {
		logs.Fatal("Invalid database driver", nil)
		return nil, errPackage.ErrUnrecognizedDriver
	}
	logs.Info("Database driver initialized successfully")

	// 2. Inicializar la conexión a la base de datos
	dbConnection := drivers.NewDatabaseConnection(driver)
	if dbConnection.Err != nil {
		logs.Fatal("Failed to connect to the database", map[string]interface{}{"error": dbConnection.Err.Error()})
		return nil, dbConnection.Err
	}
	logs.Info("Database connection initialized successfully")

	// 3. Abrir la conexión a la base de datos
	if err := dbConnection.Open(); err != nil {
		logs.Fatal("Failed to open database connection", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// 4. Iniciar migraciones solo si así está definido en la configuración
	if config.Server.RunMigration {
		err := database.RunMigrations(dbConnection.Db)
		if err != nil {
			logs.Fatal("Failed to run migrations", map[string]interface{}{"error": err.Error()})
			return nil, err
		}
	}

	return dbConnection, nil
}

// selectDatabaseDriver selecciona el driver de la base de datos según la configuración del entorno
func (app *Application) selectDatabaseDriver() drivers.DriverConfig {
	driver, ok := SupportedDrivers[config.Database.Driver]
	if !ok {
		return nil
	}
	return driver
}
