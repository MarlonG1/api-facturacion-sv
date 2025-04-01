package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/service"
	ccfInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/interfaces"
	ccfService "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/service"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/interfaces"
	contiEventPort "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/ports"
	service2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/service"
	transmissionPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	transmission "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/service"
	invalidationManager "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/interfaces"
	invalidationService "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/service"
	invoiceInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/interfaces"
	invoiceService "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/service"
	transmitter2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	batchPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/ports"
	healthPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	metricsPort "github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	testPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/cache"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/contingency"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/crypt"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/health"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/metrics"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing/signer"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/test_endpoint"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/tokens"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/batch"
	"time"
)

type ServicesContainer struct {
	repos *RepositoryContainer

	cacheManager            ports.CacheManager
	tokenManager            ports.TokenManager
	authManager             ports.AuthManager
	cryptManager            ports.CryptManager
	transmitterManager      appPorts.DTETransmitter
	haciendaAuthManager     appPorts.HaciendaAuthManager
	signerManager           appPorts.SignerManager
	dteManager              transmissionPorts.DTEManager
	sequentialManager       transmissionPorts.SequentialNumberManager
	invalidationManager     invalidationManager.InvalidationManager
	invoiceManager          invoiceInterfaces.InvoiceManager
	ccfManager              ccfInterfaces.CCFManager
	transmitterBatchManager batchPorts.BatchTransmitterPort
	contingencyEventManager contiEventPort.ContingencyEventSender
	contingencyManager      interfaces.ContingencyManager
	healthManager           healthPorts.HealthManager
	testManager             testPorts.TestManager
	metricsManager          metricsPort.MetricsManager
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

	c.tokenManager = tokens.NewJWTService(config.Server.JWTSecret, c.cacheManager)
	c.authManager = service.NewAuthService(c.tokenManager, c.repos.AuthRepo(), c.cacheManager)
	c.signerManager = signer.NewDTESigner(c.repos.AuthRepo())
	c.haciendaAuthManager = signing.NewHaciendaAuthService(c.cacheManager, c.authManager)
	c.transmitterManager = transmitter.NewMHTransmitter(c.haciendaAuthManager, c.repos.FailedSequentialNumberRepo())
	c.dteManager = transmission.NewDTEManager(c.repos.DTERepo())
	c.sequentialManager = transmission.NewSequentialNumberManager(c.repos.SequentialNumberRepo(), c.repos.AuthRepo())
	c.invoiceManager = invoiceService.NewInvoiceService(c.sequentialManager)
	c.ccfManager = ccfService.NewCCFService(c.sequentialManager)
	c.invalidationManager = invalidationService.NewInvalidationManager(c.dteManager)
	c.testManager = test_endpoint.NewTestService(c.repos.db)
	c.metricsManager = metrics.NewMetricManager(c.cacheManager)
	c.healthManager = health.NewHealthService(&health.HealthServiceConfig{
		DB:          c.repos.db,
		RedisClient: c.cacheManager.GetRedisClient(),
	})

	transmissionConf := models.NewTransmissionConfig(5*time.Second, 2*time.Minute, 2.0)
	c.transmitterBatchManager = batch.NewBatchTransmitterService(
		c.haciendaAuthManager,
		c.signerManager,
		c.repos.ContingencyRepo(),
		c.repos.FailedSequentialNumberRepo(),
		transmissionConf,
		&transmitter2.RealTimeProvider{},
		c.repos.connection,
	)

	c.contingencyEventManager = contingency.NewContingencyEventService(
		c.authManager,
		c.haciendaAuthManager,
		c.cacheManager,
		c.tokenManager,
		c.signerManager,
		c.repos.ContingencyRepo(),
		&transmitter2.RealTimeProvider{},
		c.repos.connection,
	)

	c.contingencyManager = service2.NewContingencyManager(
		c.authManager,
		c.dteManager,
		c.repos.ContingencyRepo(),
		c.haciendaAuthManager,
		c.cacheManager,
		c.tokenManager,
		c.signerManager,
		c.transmitterBatchManager,
		c.contingencyEventManager,
		&transmitter2.RealTimeProvider{},
		transmissionConf,
	)

	return nil
}

func (c *ServicesContainer) MetricsManager() metricsPort.MetricsManager {
	return c.metricsManager
}

func (c *ServicesContainer) HealthManager() healthPorts.HealthManager {
	return c.healthManager
}

func (c *ServicesContainer) TestManager() testPorts.TestManager {
	return c.testManager
}

func (c *ServicesContainer) InvalidationManager() invalidationManager.InvalidationManager {
	return c.invalidationManager
}

func (c *ServicesContainer) ContingencyManager() interfaces.ContingencyManager {
	return c.contingencyManager
}

func (c *ServicesContainer) DTEManager() transmissionPorts.DTEManager {
	return c.dteManager
}

func (c *ServicesContainer) CCFService() ccfInterfaces.CCFManager {
	return c.ccfManager
}

func (c *ServicesContainer) InvoiceService() invoiceInterfaces.InvoiceManager {
	return c.invoiceManager
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
