package config

import (
	"errors"
	"fmt"
	errPackage "github.com/MarlonG1/api-facturacion-sv/config/error"
	"github.com/spf13/viper"
	"reflect"
	"regexp"
	"strings"
)

// ============================= GUÍA DE USO =============================
// Para agregar una nueva variable de entorno, siga los siguientes pasos:
//
//  1. Defina la variable en el archivo `.env` con un nombre descriptivo.
//     Ejemplo:
//     NUEVA_VARIABLE=valor
//
//  2. Agregue un campo en la estructura `envConfig` dentro del archivo `env_structs` con la etiqueta `map-structure`.
//     Esto asegurará que se pueda mapear correctamente la variable.
//     Ejemplo:
//     type EnvConfig struct {
//     NuevaVariable string `map-structure:"NUEVA_VARIABLE"`
//     }
//
//  3. Si necesita que su nueva o nuevas validadas sigan un patrón específico, defina una constante con el patrón
//     en este archivo. Luego, en la función `ValidateConfig`, agregue una validación para el campo correspondiente.
//     Ejemplo:
//     var (
//     ExamplePattern = "^[a-zA-Z0-9_]*$"
//     )
//     ...
//     if !matchPattern(ExamplePattern, EnvConfig.NuevaVariable) {
//         return fmt.Errorf("NUEVA_VARIABLE must match the pattern")
//     }
//
//  4. Para acceder a la variable en el código, use:
//     env.NuevaVariable
//     O si fue una estructura anidada:
//     env.EstructuraAnidada.NuevaVariable
//
//     Nota: Los tipos de datos soportados son `string`, `bool`, `int`, `uint`, `float32` y `float64`.
//     El código se encargará de buscar la variable de entorno con el nombre definido en la etiqueta
//     `map-structure` y asignarla al campo correspondiente en la estructura `envConfig`.
//
// Siguiendo esos pasos, podrá agregar nuevas variables de entorno de forma sencilla y estructurada.
// ======================================================================

var (
	// URLPattern es un regex para validar URLs que pueden incluir "http" o "https",
	// permitir "localhost", nombres de dominio con TLD válidos y un puerto opcional.
	// Casos válidos:
	//   - http://example.com
	//   - https://www.example.org
	//   - http://sub.dominio.net:8080
	//   - https://localhost:3000/path
	//   - http://127.0.0.1:5000/api
	//   - http://miweb.com:65535/servicio
	// Casos inválidos:
	//   - http://example.com: (No debería permitir solo ":" sin puerto)
	//   - https://site:70000 (70000 está fuera del rango de puertos)
	//   - http://localhost:abc/ (Solo números permitidos en el puerto)
	URLPattern = "https?:\\/\\/(?:localhost|(?:www\\.)?[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}|(?:\\d{1,3}\\.){3}\\d{1,3})(:\\d{1,5})?(\\/\\S*)?"

	// PortPattern es un regex para validar números de puertos en el rango de 0 a 65535.
	// Casos válidos:
	//   - 0
	//   - 22
	//   - 80
	//   - 443
	//   - 8080
	//   - 65535
	// Casos inválidos:
	//   - -1 (No existen puertos negativos)
	//   - 65536 (Fuera del rango permitido)
	//   - 123456 (Demasiado grande)
	//   - 080 (No debería aceptar ceros a la izquierda)
	//   - 23a (No debe contener letras)
	PortPattern = "\\b(0|[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])\\b"

	// HostPattern es un regex para validar hosts de forma estricta.
	// Este patrón acepta:
	//  - "localhost".
	//  - Direcciones IPv4, donde cada octeto debe estar en el rango de 0 a 255.
	//  - Nombres de dominio complejos, incluyendo dominios con múltiples subdominios
	//    (por ejemplo, dominios de AWS), asegurando un TLD de al menos 2 letras.
	//
	// Ejemplos de casos válidos:
	//   - localhost
	//   - 127.0.0.1
	//   - 192.168.1.100
	//   - 8.8.8.8
	//   - www.example.com
	//   - ec2-54-152-121-53.compute-1.amazonaws.com
	//
	// Ejemplos de casos inválidos:
	//   - 256.256.256.256   (cada octeto debe estar entre 0 y 255)
	//   - example            (falta el TLD)
	//   - -1.2.3.4           (no se permiten números negativos)
	//   - 192.168.1          (falta un octeto)
	//   - 192.168.1.1.1      (demasiados octetos)
	HostPattern = "^(localhost|((25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)\\.){3}(25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)|((?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,}))$"
)

var (
	// AvailableDatabaseDrivers contiene los drivers de base de datos soportados.
	// Se usa para validar que el driver especificado en la configuración sea válido.
	// Si se agrega un nuevo driver, se debe agregar aquí.
	AvailableDatabaseDrivers = map[string]bool{
		"mysql":    true,
		"postgres": true,
	}
)

var EnvConfig *envConfig
var Server *server
var Database *database
var Redis *redis
var Log *log
var Signer *signer
var MHPaths *mhPaths

// InitEnvConfig inicializa la configuración del archivo .env
func InitEnvConfig(rootPath string) error {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(rootPath)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return errPackage.ErrEnvFileNotFound
		}
		return err
	}

	v.AutomaticEnv()
	EnvConfig = &envConfig{}

	// Mapear automáticamente las variables de entorno usando reflection y etiquetas
	if err := autoMapEnvKeys(v, reflect.ValueOf(EnvConfig).Elem()); err != nil {
		return err
	}

	if err := v.Unmarshal(&EnvConfig); err != nil {
		return errPackage.ErrFailedToLoadEnv
	}

	if err := ValidateConfig(); err != nil {
		return err
	}

	// Asignar las estructuras a las variables globales
	Server = &EnvConfig.Server
	Database = &EnvConfig.Database
	Redis = &EnvConfig.Redis
	Log = &EnvConfig.Log
	Signer = &EnvConfig.Signer
	MHPaths = &EnvConfig.MHPaths

	return nil
}

// ValidateConfig valida cada campo de la configuración del archivo .env
func ValidateConfig() error {
	if err := validateServerFields(); err != nil {
		return err
	}

	if err := validateDatabaseFields(); err != nil {
		return err
	}

	if err := validateRedisFields(); err != nil {
		return err
	}

	if err := validateLogFields(); err != nil {
		return err
	}

	if err := validateMHConfigFields(); err != nil {
		return err
	}

	if err := validateSignerFields(); err != nil {
		return err
	}

	return nil
}

// validateServerFields valida los campos de la estructura Server
func validateServerFields() error {
	bt := map[string]bool{
		"DEBUG":            true,
		"FORCECONTINGENCY": true,
		"RUNMIGRATION":     true,
	}
	v := reflect.ValueOf(EnvConfig.Server)

	if err := validateEnvVariables(v, bt, nil); err != nil {
		return err
	}

	if !matchPattern(PortPattern, EnvConfig.Server.Port) {
		return fmt.Errorf("SERVER_PORT must be a valid port")
	}

	if EnvConfig.Server.MaxBatchSize <= 0 || EnvConfig.Server.MaxBatchSize > 100 {
		return fmt.Errorf("MH_MAX_BATCH_SIZE must be between 1 and 100")
	}

	return nil
}

// validateDatabaseFields valida los campos de la estructura Database
func validateDatabaseFields() error {
	v := reflect.ValueOf(EnvConfig.Database)

	if err := validateEnvVariables(v, nil, nil); err != nil {
		return err
	}

	if !matchPattern(HostPattern, EnvConfig.Database.Host) {
		return fmt.Errorf("DATABASE_HOST must be a valid host")
	}

	if !matchPattern(PortPattern, EnvConfig.Database.Port) {
		return fmt.Errorf("DATABASE_PORT must be a valid port")
	}

	if !AvailableDatabaseDrivers[EnvConfig.Database.Driver] {
		return fmt.Errorf("DATABASE_DRIVER must be a valid driver")
	}

	return nil
}

// validateRedisFields valida los campos de la estructura Redis
func validateRedisFields() error {
	v := reflect.ValueOf(EnvConfig.Redis)
	ex := []string{strings.ToUpper("PASSWORD")}

	if err := validateEnvVariables(v, nil, ex); err != nil {
		return err
	}

	if !matchPattern(HostPattern, EnvConfig.Redis.Host) {
		return fmt.Errorf("REDIS_HOST must be a valid host")
	}

	if !matchPattern(PortPattern, EnvConfig.Redis.Port) {
		return fmt.Errorf("REDIS_PORT must be a valid port")
	}

	return nil
}

// validateLogFields valida los campos de la estructura Log
func validateLogFields() error {
	bt := map[string]bool{
		"FILELOGGING": true,
	}
	v := reflect.ValueOf(EnvConfig.Log)

	if err := validateEnvVariables(v, bt, nil); err != nil {
		return err
	}

	return nil
}

// validateMHConfigFields valida los campos de la estructura MHPaths
func validateMHConfigFields() error {
	v := reflect.ValueOf(EnvConfig.MHPaths)

	if err := validateEnvVariables(v, nil, nil); err != nil {
		return err
	}

	for i := 0; i < v.NumField(); i++ {
		t := v.Type()
		f := v.Field(i)

		if !matchPattern(URLPattern, f.String()) {
			return fmt.Errorf("%s must be a valid URL", strings.ToUpper(t.Field(i).Name))
		}
	}

	return nil
}

// validateSignerFields valida los campos de la estructura Signer
func validateSignerFields() error {
	v := reflect.ValueOf(EnvConfig.Signer)

	if err := validateEnvVariables(v, nil, nil); err != nil {
		return err
	}

	for i := 0; i < v.NumField(); i++ {
		t := v.Type()
		f := v.Field(i)

		if !matchPattern(URLPattern, f.String()) {
			return fmt.Errorf("%s must be a valid URL", strings.ToUpper(t.Field(i).Name))
		}
	}

	return nil
}

// validateEnvVariables valida que los campos de la estructura sean requeridos y del tipo correcto
func validateEnvVariables(v reflect.Value, bt map[string]bool, exceptions []string) error {
	t := v.Type()
	// Se crea un mapa con las excepciones
	exMap := make(map[string]bool)
	for _, ex := range exceptions {
		exMap[strings.ToUpper(ex)] = true
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fn := strings.ToUpper(t.Field(i).Name)

		// Solo se valida si el campo es de tipo boolean
		if bt != nil && bt[fn] {
			if f.Kind() != reflect.Bool {
				return fmt.Errorf("%s must be boolean", fn)
			}
			continue
		}

		// Comprobar el tipo de campo y validarlo adecuadamente
		switch f.Kind() {
		case reflect.String:
			// Si el campo está en la lista de excepciones, no se valida
			if exMap[fn] {
				continue
			}

			// Si el campo está vacío, regresa un error
			if f.String() == "" {
				return fmt.Errorf("%s is required", fn)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if exMap[fn] {
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if exMap[fn] {
				continue
			}
		case reflect.Float32, reflect.Float64:
			if exMap[fn] {
				continue
			}
		default:
			return fmt.Errorf("unsupported type %s for field %s", f.Kind(), fn)
		}
	}

	return nil
}

// autoMapEnvKeys mapea las variables de entorno a las llaves de la estructura
func autoMapEnvKeys(v *viper.Viper, val reflect.Value) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		t := typ.Field(i)

		// Si el campo es una estructura, se llama recursivamente
		if f.Kind() == reflect.Struct {
			if err := autoMapEnvKeys(v, f); err != nil {
				return err
			}
			continue
		}

		// Obtener etiqueta de la estructura
		envVar := t.Tag.Get("map-structure")
		if envVar == "" {
			continue
		}

		// Asignar el valor segun el tipo de dato
		switch f.Kind() {
		case reflect.String:
			f.SetString(v.GetString(envVar))
		case reflect.Bool:
			f.SetBool(v.GetBool(envVar))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(v.GetInt64(envVar))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f.SetUint(v.GetUint64(envVar))
		case reflect.Float32, reflect.Float64:
			f.SetFloat(v.GetFloat64(envVar))
		default:
			// Tipos no soportados
			return fmt.Errorf("unsupported type %s", f.Kind())
		}
	}

	return nil
}

// matchPattern compara un patrón con un valor
func matchPattern(pattern, value string) bool {
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}
