package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type RedisTokenCache struct {
	client       *redis.Client
	cryptService ports.CryptManager
	ctx          context.Context
}

// NewRedisTokenCache crea una nueva instancia de RedisTokenCache
func NewRedisTokenCache(config *config.RedisConfig, cryptService ports.CryptManager) (ports.CacheManager, error) {
	opt, err := redis.ParseURL(config.GetURL())
	if err != nil {
		logs.Error("Failed to parse Redis URL", map[string]interface{}{
			"url":   config.GetURL(),
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"NewRedisTokenCache",
			"failed to parse Redis URL",
			err,
		)
	}

	client := redis.NewClient(opt)

	// Verificar conexión
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		logs.Error("Failed to connect to Redis", map[string]interface{}{
			"host":  config.Host,
			"port":  config.Port,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"NewRedisTokenCache",
			"failed to connect to Redis",
			err,
		)
	}

	logs.Info("Successfully connected to Redis", map[string]interface{}{
		"host": config.Host,
		"port": config.Port,
	})
	return &RedisTokenCache{
		client:       client,
		cryptService: cryptService,
		ctx:          ctx,
	}, nil
}

// Set guarda un token en Redis con un tiempo de vida determinado
func (c *RedisTokenCache) Set(key string, saveInfo []byte, ttl time.Duration) error {
	err := c.client.Set(c.ctx, key, saveInfo, ttl).Err()
	if err != nil {
		logs.Error("Failed to set value in Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"Set",
			"failed to set value in Redis",
			err,
		)
	}

	logs.Info("Value set successfully in Redis", map[string]interface{}{
		"key": key,
		"ttl": ttl.Seconds(),
	})
	return nil
}

// SetCredentials guarda las credenciales de Hacienda en Redis con un tiempo de vida determinado
func (c *RedisTokenCache) SetCredentials(token string, creds *models.HaciendaCredentials, ttl time.Duration) error {
	key := fmt.Sprintf("hacienda:credentials:%s", token)

	encryptStruct, err := c.cryptService.EncryptStruct(token, *creds)
	if err != nil {
		logs.Error("Failed to encrypt credentials", map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"SetCredentials",
			"failed to encrypt credentials",
			err,
		)
	}

	err = c.client.Set(c.ctx, key, encryptStruct, ttl).Err()
	if err != nil {
		logs.Error("Failed to set credentials in Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"SetCredentials",
			"failed to set credentials in Redis",
			err,
		)
	}

	logs.Info("Credentials set successfully in Redis", map[string]interface{}{
		"key": key,
		"ttl": ttl.Seconds(),
	})
	return nil
}

// Get obtiene un token de Redis y lo convierte en un AuthClaims
func (c *RedisTokenCache) Get(key string) (string, error) {
	cacheInfo, err := c.client.Get(c.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		logs.Error("Token not found in Redis", map[string]interface{}{
			"key": key,
		})
		return "", shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"Get",
			"token not found in Redis",
			errPackage.ErrTokenNotFound,
		)
	}
	if err != nil {
		logs.Error("Failed to get value from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return "", shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"Get",
			"failed to get value from Redis",
			err,
		)
	}

	logs.Info("Value retrieved successfully from Redis")
	return cacheInfo, nil
}

// GetCredentials obtiene las credenciales de Hacienda de Redis y las convierte en un HaciendaCredentials
func (c *RedisTokenCache) GetCredentials(token string) (*models.HaciendaCredentials, error) {
	key := fmt.Sprintf("hacienda:credentials:%s", token)

	cipherStruct, err := c.client.Get(c.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		logs.Error("Credentials not found in Redis", map[string]interface{}{
			"key": key,
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"GetCredentials",
			"credentials not found in Redis",
			errPackage.ErrTokenNotFound,
		)
	}
	if err != nil {
		logs.Error("Failed to get credentials from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"GetCredentials",
			"failed to get credentials from Redis",
			err,
		)
	}

	creds, err := c.cryptService.DecryptStruct(token, cipherStruct)
	if err != nil {
		logs.Error("Failed to decrypt credentials", map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"GetCredentials",
			"failed to decrypt credentials",
			err,
		)
	}

	logs.Info("Credentials decrypted successfully", map[string]interface{}{
		"key": key,
	})
	return &creds, nil
}

// Delete elimina un token de Redis
func (c *RedisTokenCache) Delete(token string) error {
	key := "token:" + token
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		logs.Error("Failed to delete token from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"Delete",
			"failed to delete token from Redis",
			err,
		)
	}

	logs.Info("Token deleted successfully from Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

// Close cierra la conexión con Redis
func (c *RedisTokenCache) Close() error {
	err := c.client.Close()
	if err != nil {
		logs.Error("Failed to close Redis client", map[string]interface{}{
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"Close",
			"failed to close Redis client",
			err,
		)
	}

	logs.Info("Redis client closed successfully", map[string]interface{}{})
	return nil
}

func (c *RedisTokenCache) RPush(key string, value []byte) error {
	err := c.client.RPush(c.ctx, key, value).Err()
	if err != nil {
		logs.Error("Failed to RPush value to Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"RPush",
			"failed to RPush value to Redis",
			err,
		)
	}

	logs.Info("Value RPushed successfully to Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

func (c *RedisTokenCache) LPush(key string, value []byte) error {
	err := c.client.LPush(c.ctx, key, value).Err()
	if err != nil {
		logs.Error("Failed to LPush value to Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"LPush",
			"failed to LPush value to Redis",
			err,
		)
	}

	logs.Info("Value LPushed successfully to Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

func (c *RedisTokenCache) LRange(key string, start, stop int64) ([]string, error) {
	result, err := c.client.LRange(c.ctx, key, start, stop).Result()
	if err != nil {
		logs.Error("Failed to get range from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"LRange",
			"failed to get range from Redis",
			err,
		)
	}

	return result, nil
}

func (c *RedisTokenCache) LLen(key string) (int64, error) {
	length, err := c.client.LLen(c.ctx, key).Result()
	if err != nil {
		logs.Error("Failed to get list length from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"LLen",
			"failed to get list length from Redis",
			err,
		)
	}

	return length, nil
}

func (c *RedisTokenCache) LTrim(key string, start, stop int64) error {
	err := c.client.LTrim(c.ctx, key, start, stop).Err()
	if err != nil {
		logs.Error("Failed to trim list in Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"RedisTokenCache",
			"LTrim",
			"failed to trim list in Redis",
			err,
		)
	}

	logs.Info("List trimmed successfully in Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}
