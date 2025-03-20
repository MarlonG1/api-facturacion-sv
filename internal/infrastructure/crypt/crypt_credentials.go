package crypt

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gtank/cryptopasta"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// EncryptStruct encripta una estructura de HaciendaCredentials y la convierte en un string
func EncryptStruct(token string, data models.HaciendaCredentials) (string, error) {
	key := DeriveKeyFromToken(token)

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
func DecryptStruct(token string, data string) (models.HaciendaCredentials, error) {
	var creds models.HaciendaCredentials
	key := DeriveKeyFromToken(token)

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
