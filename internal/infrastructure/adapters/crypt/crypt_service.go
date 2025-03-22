package crypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/gtank/cryptopasta"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type CryptService struct{}

func NewCryptService() ports.CryptManager {
	return &CryptService{}
}

// GenerateAPIKey genera una API key segura
func (cs *CryptService) GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "GenerateAPIKey", "error generating random bytes", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPISecret genera un API secret seguro
func (cs *CryptService) GenerateAPISecret() (string, error) {
	bytes := make([]byte, 48)
	if _, err := rand.Read(bytes); err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "GenerateAPISecret", "error generating random bytes", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// DeriveKeyFromToken deriva una clave de un token dado a través de SHA-256
func (cs *CryptService) deriveKeyFromToken(token string) *[32]byte {
	hash := sha256.Sum256([]byte(token))
	return &hash
}

// EncryptStruct encripta una estructura de HaciendaCredentials y la convierte en un string
func (cs *CryptService) EncryptStruct(token string, data models.HaciendaCredentials) (string, error) {
	key := cs.deriveKeyFromToken(token)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "EncryptStruct", "error marshalling data", err)
	}

	secret, err := cryptopasta.Encrypt(jsonData, key)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "EncryptStruct", "error encrypting data", err)
	}

	return base64.StdEncoding.EncodeToString(secret), nil
}

// DecryptStruct desencripta un string y lo convierte en una estructura de HaciendaCredentials
func (cs *CryptService) DecryptStruct(token string, data string) (models.HaciendaCredentials, error) {
	var creds models.HaciendaCredentials
	key := cs.deriveKeyFromToken(token)

	preDecodeData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return models.HaciendaCredentials{}, shared_error.NewGeneralServiceError("Utils", "DecryptStruct", "error decoding base64 data", err)
	}

	decrypted, err := cryptopasta.Decrypt(preDecodeData, key)
	if err != nil {
		return models.HaciendaCredentials{}, shared_error.NewGeneralServiceError("Utils", "DecryptStruct", "error decrypting data", err)
	}

	err = json.Unmarshal(decrypted, &creds)
	if err != nil {
		return models.HaciendaCredentials{}, shared_error.NewGeneralServiceError("Utils", "DecryptStruct", "error unmarshalling data", err)
	}

	return creds, nil
}

// GenerateBulkAPIKeys es una funcion de tipo bulk que genera una cantidad determinada de API KEYS y API SECRETs
func (cs *CryptService) GenerateBulkAPIKeys(amount int) ([]string, []string, error) {
	var err error
	keys := make([]string, 0)
	secrets := make([]string, 0)

	for i := 0; i < amount; i++ {
		keys[i], err = cs.GenerateAPIKey()
		if err != nil {
			return nil, nil, err
		}
		secrets[i], err = cs.GenerateAPISecret()
		if err != nil {
			return nil, nil, err
		}
	}

	return keys, secrets, nil
}
