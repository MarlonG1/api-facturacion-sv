package main

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/cmd/setup"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/server"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"

	errPackage "github.com/MarlonG1/api-facturacion-sv/config/error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

var (
	SupportedDrivers = map[string]drivers.DriverConfig{
		"mysql":    drivers.NewMysqlDriver(),
		"postgres": drivers.NewPostgresDriver(),
	}
)

// main es la función principal de la aplicación que se encarga de inicializar los componentes necesarios para el
// correcto funcionamiento de la aplicación.
func main() {
	// 0. Obtener el root path del proyecto
	rootPath := utils.FindProjectRoot()

	// 1. Inicializar la configuración del entorno
	err := config.InitEnvConfig(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Inicializar el logger
	err = logs.InitLogger(config.Log.Level, config.Log.Path)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}
	logs.Info("Logger initialized successfully")

	// 3. Inicializar el tiempo global
	err = utils.TimeInit()
	if err != nil {
		logs.Fatal("Failed to initialize global time", map[string]interface{}{"error": err.Error()})
		return
	}

	// 3. Iniciar la configuración de la base de datos y las migraciones
	dbConnection, err := initDatabaseConfigurations()
	if err != nil {
		logs.Fatal("Failed to initialize database configurations", map[string]interface{}{"error": err.Error()})
		return
	}

	// 4. Inicializar el contenedor de dependencias
	container := bootstrap.NewContainer(dbConnection.Db)
	err = container.Initialize()
	if err != nil {
		logs.Error("Failed to initialize container", map[string]interface{}{"error": err.Error()})
		return
	}

	// 5. Inicializar el servidor
	sv := server.Initialize(container)

	// 6. Inicializar los jobs
	err = setup.SetupJobs(container.Services().ContingencyManager(), config.Server.AmbientCode)
	if err != nil {
		logs.Error("Failed to setup jobs", map[string]interface{}{"error": err.Error()})
		return
	}

	// 7. Iniciar el servidor
	logs.Info("Server started successfully", map[string]interface{}{"port": config.Server.Port})
	if err = sv.Start(); err != nil {
		logs.Fatal("Failed to start server", map[string]interface{}{"error": err.Error()})
		return
	}

}

func initDatabaseConfigurations() (*drivers.DbConnection, error) {
	// 1. Seleccionar el driver de la base de datos
	driver := selectDatabaseDriver()
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

	// 4. Iniciar migraciones
	err := database.RunMigrations(dbConnection.Db)
	if err != nil {
		logs.Fatal("Failed to run migrations", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return dbConnection, nil
}

// selectDatabaseDriver selecciona el driver de la base de datos según la configuración del entorno.
// Retorna una instancia de la interfaz DriverConfig.
// Si ha agregado un nuevo driver de base de datos, debe agregar un nuevo case en el switch.
func selectDatabaseDriver() drivers.DriverConfig {
	driver, ok := SupportedDrivers[config.Database.Driver]
	if !ok {
		return nil
	}
	return driver
}
