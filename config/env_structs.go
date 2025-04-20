package config

// envConfig es una estructura que contiene la configuración del archivo .env
type envConfig struct {
	Server   server
	Database database
	Redis    redis
	Log      log
	Signer   signer
	MHPaths  mhPaths
}

// server es una estructura que contiene la configuración del servidor
type server struct {
	Port             string `map-structure:"SERVER_PORT"`
	MaxBatchSize     int    `map-structure:"MH_MAX_BATCH_SIZE"`
	JWTSecret        string `map-structure:"JWT_SECRET"`
	AmbientCode      string `map-structure:"MH_AMBIENT_CODE"`
	Debug            bool   `map-structure:"DEBUG"`
	RunMigration     bool   `map-structure:"RUN_MIGRATION"`
	AdminEmail       string `map-structure:"ADMIN_EMAIL"`
	ForceContingency bool   `map-structure:"FORCE_CONTINGENCY"`
	AppLang          string `map-structure:"APP_LANG"`
}

// database es una estructura que contiene la configuración de la base de datos
type database struct {
	Host     string `map-structure:"DB_HOST"`
	Port     string `map-structure:"DB_PORT"`
	Name     string `map-structure:"DB_DATABASE"`
	User     string `map-structure:"DB_USERNAME"`
	Password string `map-structure:"DB_PASSWORD"`
	Charset  string `map-structure:"DB_CHARSET"`
	Driver   string `map-structure:"DB_DRIVER"`
}

// redis es una estructura que contiene la configuración de redis
type redis struct {
	Host     string `map-structure:"REDIS_HOST"`
	Port     string `map-structure:"REDIS_PORT"`
	Password string `map-structure:"REDIS_PASSWORD"`
}

// log es una estructura que contiene la configuración de los logs
type log struct {
	Level       string `map-structure:"LOG_LEVEL"`
	Path        string `map-structure:"LOG_PATH"`
	FileLogging bool   `map-structure:"LOG_FILE_LOGGING"`
}

// signer es una estructura que contiene la configuración del firmante
type signer struct {
	Path   string `map-structure:"SIGNER_PATH"`
	Health string `map-structure:"SIGNER_HEALTH"`
}

// mhPaths es una estructura que contiene las rutas de los servicios de MH
type mhPaths struct {
	AuthURL                 string `map-structure:"MH_AUTH_URL"`
	ReceptionURL            string `map-structure:"MH_RECEPTION_URL"`
	LoteReceptionURL        string `map-structure:"MH_LOTE_RECEPTION_URL"`
	ReceptionConsultURL     string `map-structure:"MH_RECEPTION_CONSULT_URL"`
	LoteReceptionConsultURL string `map-structure:"MH_RECEPTION_CONSULT_LOTE_URL"`
	ContingencyURL          string `map-structure:"MH_CONTINGENCY_URL"`
	NullifyURL              string `map-structure:"MH_NULLIFY_URL"`
}
