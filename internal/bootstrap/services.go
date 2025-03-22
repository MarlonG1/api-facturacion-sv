package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/service"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/cache"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/crypt"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/tokens"
)

type ServicesContainer struct {
	repos *RepositoryContainer

	cacheManager ports.CacheManager
	tokenManager ports.TokenManager
	authManager  ports.AuthManager
	cryptManager ports.CryptManager
}

func NewServicesContainer(repos *RepositoryContainer) *ServicesContainer {
	return &ServicesContainer{
		repos: repos,
	}
}

func (c *ServicesContainer) Initialize() error {
	var err error

	c.cryptManager = crypt.NewCryptService()
	c.cacheManager, err = cache.NewRedisTokenCache(config.NewRedisConfig(), c.cryptManager)
	if err != nil {
		return err
	}

	c.tokenManager = tokens.NewJWTService(env.Server.JWTSecret, c.cacheManager)
	c.authManager = service.NewAuthService(c.tokenManager, c.repos.AuthRepo(), c.cacheManager)

	return nil
}

func (c *ServicesContainer) CacheManager() ports.CacheManager {
	return c.cacheManager
}

func (c *ServicesContainer) TokenManager() ports.TokenManager {
	return c.tokenManager
}

func (c *ServicesContainer) AuthManager() ports.AuthManager {
	return c.authManager
}

func (c *ServicesContainer) CryptManager() ports.CryptManager {
	return c.cryptManager
}
