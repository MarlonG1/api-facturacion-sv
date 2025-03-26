package ports

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/go-redis/redis/v8"
	"time"
)

// CacheManager interface para abstracción del cache
type CacheManager interface {
	Set(key string, claims []byte, ttl time.Duration) error                                       // Set guarda un token en el cache
	SetCredentials(token string, cipherInfo *models.HaciendaCredentials, ttl time.Duration) error // SetCredentials guarda las credenciales en el cache
	GetCredentials(token string) (*models.HaciendaCredentials, error)                             // GetCredentials obtiene las credenciales del cache
	Get(key string) (string, error)                                                               // Get obtiene un token del cache
	Delete(token string) error                                                                    // Delete elimina un token del cache
	GetRedisClient() *redis.Client                                                                // GetRedisClient retorna el cliente de Redis
	CacheListManager
}

// CacheListManager define el comportamiento para la gestión de listas en caché
type CacheListManager interface {
	RPush(key string, value []byte) error                   // Añade elemento al final de la lista
	LPush(key string, value []byte) error                   // Añade elemento al inicio de la lista
	LRange(key string, start, stop int64) ([]string, error) // Obtiene rango de elementos de la lista
	LLen(key string) (int64, error)                         // Obtiene longitud de la lista
	LTrim(key string, start, stop int64) error              // Mantiene solo el rango especificado
}

// TokenManager define el comportamiento para la gestión de tokens
type TokenManager interface {
	GenerateToken(claims *models.AuthClaims) (string, error)                                     // GenerateToken genera un nuevo token JWT con los claims proporcionados
	ValidateToken(token string) (*models.AuthClaims, error)                                      // ValidateToken valida un token y retorna sus claims
	RevokeToken(token string) error                                                              // RevokeToken revoca un token específico
	SaveTimestampsForContingency(issuedAt, expiresAt time.Time, claims *models.AuthClaims) error // SaveTimestampsForContingency guarda los timestamps de un token en contingencia
	GetSecretKey() string                                                                        // GetSecretKey retorna la clave secreta para firmar los tokens
}
