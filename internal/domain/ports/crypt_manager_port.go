package ports

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"

// CryptManager determina el comportamiento de un manager de encriptación de datos
type CryptManager interface {
	// GenerateAPIKey genera un API Key aleatorio
	GenerateAPIKey() (string, error)
	// GenerateAPISecret genera un API Secret aleatorio
	GenerateAPISecret() (string, error)
	// EncryptStruct encripta una estructura de HaciendaCredentials y la convierte en un string
	EncryptStruct(token string, data models.HaciendaCredentials) (string, error)
	// DecryptStruct desencripta un string y lo convierte en una estructura de HaciendaCredentials
	DecryptStruct(token string, data string) (models.HaciendaCredentials, error)
	// GenerateBulkAPIKeys genera una cantidad de API Keys aleatorios
	GenerateBulkAPIKeys(amount int) ([]string, []string, error)
}
