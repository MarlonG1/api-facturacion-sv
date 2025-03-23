package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/service"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/interfaces"
	dteServices "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/service"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/cache"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/crypt"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing/signer"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/tokens"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter"
)

type ServicesContainer struct {
	repos *RepositoryContainer

	cacheManager        ports.CacheManager
	tokenManager        ports.TokenManager
	authManager         ports.AuthManager
	cryptManager        ports.CryptManager
	transmitterManager  appPorts.DTETransmitter
	haciendaAuthManager appPorts.HaciendaAuthManager
	signerManager       appPorts.SignerManager
	invoiceService      interfaces.InvoiceManager
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
	c.signerManager = signer.NewDTESigner(c.repos.AuthRepo())
	c.transmitterManager = transmitter.NewMHTransmitter(c.haciendaAuthManager)
	c.haciendaAuthManager = signing.NewHaciendaAuthService(c.cacheManager, c.authManager)
	c.invoiceService = dteServices.NewInvoiceService(c.repos.SequentialNumberRepo())

	return nil
}

func (c *ServicesContainer) InvoiceService() interfaces.InvoiceManager {
	return c.invoiceService
}

func (c *ServicesContainer) TransmitterManager() appPorts.DTETransmitter {
	return c.transmitterManager
}

func (c *ServicesContainer) SignerManager() appPorts.SignerManager {
	return c.signerManager
}

func (c *ServicesContainer) HaciendaAuthManager() appPorts.HaciendaAuthManager {
	return c.haciendaAuthManager
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
