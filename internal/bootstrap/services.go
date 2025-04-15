package bootstrap

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/service/strategies"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/cache"
	adapterContingecy "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/contingency"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/crypt"
	adapterHealth "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/health"
	adapterMetric "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/metrics"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/signing/signer"
	adapterTest "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/test_endpoint"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/tokens"
	adapterTransmitter "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter"
	batch "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/batch"
)

type ServicesContainer struct {
	repos *RepositoryContainer

	cacheManager            ports.CacheManager
	tokenManager            ports.TokenManager
	authManager             auth.AuthManager
	cryptManager            ports.CryptManager
	transmitterManager      appPorts.DTETransmitter
	haciendaAuthManager     appPorts.HaciendaAuthManager
	signerManager           appPorts.SignerManager
	dteManager              dte_documents.DTEManager
	sequentialManager       dte_documents.SequentialNumberManager
	invalidationManager     invalidation.InvalidationManager
	invoiceManager          invoice.InvoiceManager
	ccfManager              ccf.CCFManager
	transmitterBatchManager transmitter.BatchTransmitterPort
	contingencyEventManager contingency.ContingencyEventSender
	contingencyManager      contingency.ContingencyManager
	healthManager           health.HealthManager
	testManager             test_endpoint.TestManager
	metricsManager          metrics.MetricsManager
	retentionManager        retention.RetentionManager
	creditNoteManager       credit_note.CreditNoteManager
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
	c.authManager = strategies.NewAuthService(c.tokenManager, c.repos.AuthRepo(), c.cacheManager)
	c.signerManager = signer.NewDTESigner(c.repos.AuthRepo())
	c.haciendaAuthManager = signing.NewHaciendaAuthService(c.cacheManager, c.authManager)
	c.transmitterManager = adapterTransmitter.NewMHTransmitter(c.haciendaAuthManager, c.repos.FailedSequentialNumberRepo())
	c.dteManager = dte_documents.NewDTEService(c.repos.DTERepo())
	c.sequentialManager = dte_documents.NewSequentialNumberService(c.repos.SequentialNumberRepo(), c.repos.AuthRepo())
	c.invoiceManager = invoice.NewInvoiceService(c.sequentialManager)
	c.ccfManager = ccf.NewCCFService(c.sequentialManager)
	c.invalidationManager = invalidation.NewInvalidationService(c.dteManager)
	c.retentionManager = retention.NewRetentionService(c.sequentialManager, c.dteManager)
	c.creditNoteManager = credit_note.NewCreditNoteService(c.sequentialManager, c.dteManager)
	c.testManager = adapterTest.NewTestService(c.repos.db)
	c.metricsManager = adapterMetric.NewMetricService(c.cacheManager)
	c.healthManager = adapterHealth.NewHealthService(&adapterHealth.HealthServiceConfig{
		DB: c.repos.db,
	})

	transmissionConf := models.NewTransmissionConfig(5*time.Second, 2*time.Minute, 2.0)
	c.transmitterBatchManager = batch.NewBatchTransmitterService(
		c.haciendaAuthManager,
		c.signerManager,
		c.repos.ContingencyRepo(),
		transmissionConf,
		&transmitter.RealTimeProvider{},
		c.repos.connection,
	)

	c.contingencyEventManager = adapterContingecy.NewContingencyEventService(
		c.authManager,
		c.haciendaAuthManager,
		c.cacheManager,
		c.tokenManager,
		c.signerManager,
		c.repos.ContingencyRepo(),
		&transmitter.RealTimeProvider{},
		c.repos.connection,
	)

	c.contingencyManager = contingency.NewContingencyManager(
		c.authManager,
		c.dteManager,
		c.repos.ContingencyRepo(),
		c.haciendaAuthManager,
		c.cacheManager,
		c.tokenManager,
		c.signerManager,
		c.transmitterBatchManager,
		c.contingencyEventManager,
		&transmitter.RealTimeProvider{},
		transmissionConf,
	)

	return nil
}

func (c *ServicesContainer) CreditNoteManager() credit_note.CreditNoteManager {
	return c.creditNoteManager
}

func (c *ServicesContainer) RetentionManager() retention.RetentionManager {
	return c.retentionManager
}

func (c *ServicesContainer) MetricsManager() metrics.MetricsManager {
	return c.metricsManager
}

func (c *ServicesContainer) HealthManager() health.HealthManager {
	return c.healthManager
}

func (c *ServicesContainer) TestManager() test_endpoint.TestManager {
	return c.testManager
}

func (c *ServicesContainer) InvalidationManager() invalidation.InvalidationManager {
	return c.invalidationManager
}

func (c *ServicesContainer) ContingencyManager() contingency.ContingencyManager {
	return c.contingencyManager
}

func (c *ServicesContainer) DTEManager() dte_documents.DTEManager {
	return c.dteManager
}

func (c *ServicesContainer) CCFService() ccf.CCFManager {
	return c.ccfManager
}

func (c *ServicesContainer) InvoiceService() invoice.InvoiceManager {
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

func (c *ServicesContainer) AuthManager() auth.AuthManager {
	return c.authManager
}

func (c *ServicesContainer) CryptManager() ports.CryptManager {
	return c.cryptManager
}
