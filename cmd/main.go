package main

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config/database_drivers"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	errPackage "github.com/MarlonG1/api-facturacion-sv/config/error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

// main es la función principal de la aplicación que se encarga de inicializar los componentes necesarios para el
// correcto funcionamiento de la aplicación.
func main() {
	// 1. Inicializar la configuración del entorno
	err := env.InitEnvConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Inicializar el logger
	err = logs.InitLogger()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}
	logs.Info("Logger initialized successfully")

	// 3. Iniciar la configuración de la base de datos y las migraciones
	// TODO: Continuar con las la inicializacion de contenedores, dbConnection omitido por simplicidad
	_, err = initDatabaseConfigurations()
	if err != nil {
		logs.Fatal("Failed to initialize database configurations", map[string]interface{}{"error": err.Error()})
		return
	}

}

func initDatabaseConfigurations() (*database_drivers.DbConnection, error) {
	// 1. Seleccionar el driver de la base de datos
	driver := selectDatabaseDriver()
	if driver == nil {
		logs.Fatal("Invalid database driver", nil)
		return nil, errPackage.ErrUnrecognizedDriver
	}
	logs.Info("Database driver initialized successfully")

	// 2. Inicializar la conexión a la base de datos
	dbConnection := database_drivers.NewDatabaseConnection(driver)
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
func selectDatabaseDriver() database_drivers.DriverConfig {
	switch env.Database.Driver {
	case "mysql":
		return database_drivers.NewMysqlDriver()
	case "postgres":
		return database_drivers.NewPostgresDriver()
	default:
		return nil
	}
}
