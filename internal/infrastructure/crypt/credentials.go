package crypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// GenerateAPIKey genera una API key segura
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "GenerateAPIKey", "error generating random bytes", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPISecret genera un API secret seguro
func GenerateAPISecret() (string, error) {
	bytes := make([]byte, 48)
	if _, err := rand.Read(bytes); err != nil {
		return "", shared_error.NewGeneralServiceError("Utils", "GenerateAPISecret", "error generating random bytes", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// DeriveKeyFromToken deriva una clave de un token dado a través de SHA-256
func DeriveKeyFromToken(token string) *[32]byte {
	hash := sha256.Sum256([]byte(token))
	return &hash
}
