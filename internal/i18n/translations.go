package i18n

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// translateMutex protege el acceso a las traducciones
	translateMutex sync.RWMutex
	// translations almacena todas las traducciones por idioma
	translations = make(map[string]map[string]string)
	// initialized indica si el sistema ya fue inicializado
	initialized = false
	// lang almacena el idioma actual
	globalLang = "en"
)

// InitTranslations carga los archivos de traducción al inicio
func InitTranslations(configPath, lang string) error {
	translateMutex.Lock()
	defer translateMutex.Unlock()

	if initialized {
		return nil
	}

	// Cargar traducciones para cada idioma soportado
	globalLang = strings.ToLower(lang)
	v := viper.New()

	// Configurar viper para el archivo de errores
	errorsFile := filepath.Join(configPath, fmt.Sprintf("%s.yaml", globalLang))
	v.SetConfigFile(errorsFile)
	v.SetConfigType("yaml")

	// Leer el archivo
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	processKeys(v, globalLang)

	initialized = true
	return nil
}

// processKeys procesa recursivamente las claves del archivo de configuración
func processKeys(v *viper.Viper, lang string) {
	// Inicializar el mapa para este idioma si no existe
	if translations[lang] == nil {
		translations[lang] = make(map[string]string)
	}

	// Obtener todas las claves del archivo
	allKeys := v.AllKeys()

	for _, key := range allKeys {
		// Obtener el valor como interfaz
		value := v.Get(key)

		// Procesar según el tipo
		switch val := value.(type) {
		case string:
			// Es una traducción directa
			translations[lang][key] = val
		case map[string]interface{}:
			// Es un grupo anidado, procesar recursivamente
			for subKey, subVal := range val {
				if strVal, ok := subVal.(string); ok {
					nestedKey := key + "." + subKey
					translations[lang][nestedKey] = strVal
				}
			}
		}
	}
}

func TranslateServiceArgs(code string, params ...interface{}) string {
	code = fmt.Sprintf("service_errors.%s", strings.ToLower(code))
	return Translate(code, params...)
}

// Translate obtiene la traducción para un código y lenguaje
func Translate(code string, params ...interface{}) string {
	code = strings.ToLower(code)
	translateMutex.RLock()
	defer translateMutex.RUnlock()

	if !initialized {
		return code
	}

	normalizedLang := normalizeLanguage(globalLang)
	template, found := translations[normalizedLang][code]
	if !found {
		// Si no se encuentra en el idioma solicitado, intentar con inglés
		template, found = translations["en"][code]
		if !found {
			// Si tampoco se encuentra en inglés, devolver un mensaje por defecto
			return fmt.Sprintf("[%s] %s", code, translations[normalizedLang]["validation_errors.unknownerror"])
		}
	}

	// Aplicar los parámetros si hay
	if len(params) > 0 && strings.Contains(template, "%") {
		return fmt.Sprintf(template, params...)
	}

	return template
}

// normalizeLanguage normaliza el código de idioma
func normalizeLanguage(lang string) string {
	// Convertir a minúsculas y tomar solo los primeros dos caracteres
	simpleLang := strings.ToLower(lang)
	if len(simpleLang) >= 2 {
		simpleLang = simpleLang[:2]
	}

	// Verificar si es un idioma soportado
	switch simpleLang {
	case "es":
		return "es"
	default:
		return "en"
	}
}

func ForceReload(configPath string) error {
	translateMutex.Lock()
	initialized = false
	translateMutex.Unlock()

	return InitTranslations(configPath, "en")
}
