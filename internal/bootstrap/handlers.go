package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
)

type HandlerContainer struct {
	useCases *UseCaseContainer
	services *ServicesContainer

	authHandler        *handlers.AuthHandler
	dteHandler         *handlers.DTEHandler
	healthHandler      *handlers.HealthHandler
	testHandler        *handlers.TestHandler
	metricsHandler     *handlers.MetricsHandler
	contingencyHandler *helpers.ContingencyHandler
}

func NewHandlerContainer(useCases *UseCaseContainer, services *ServicesContainer) *HandlerContainer {
	return &HandlerContainer{
		useCases: useCases,
		services: services,
	}
}

func (c *HandlerContainer) Initialize() {
	c.contingencyHandler = helpers.NewContingencyHandler(c.services.contingencyManager)
	c.healthHandler = handlers.NewHealthHandler(c.services.HealthManager())
	c.testHandler = handlers.NewTestHandler(c.services.TestManager())
	c.authHandler = handlers.NewAuthHandler(c.useCases.AuthUseCase())
	c.metricsHandler = handlers.NewMetricsHandler(c.services.MetricsManager())
	c.dteHandler = handlers.NewDTEHandler(c.useCases.InvoiceUseCase(), c.useCases.CCFUseCase(), c.useCases.DTEConsultUseCase(), c.useCases.InvalidationUseCase(), c.contingencyHandler)
}

func (c *HandlerContainer) MetricsHandler() *handlers.MetricsHandler {
	return c.metricsHandler
}

func (c *HandlerContainer) HealthHandler() *handlers.HealthHandler {
	return c.healthHandler
}

func (c *HandlerContainer) TestHandler() *handlers.TestHandler {
	return c.testHandler
}

func (c *HandlerContainer) DTEHandler() *handlers.DTEHandler {
	return c.dteHandler
}

func (c *HandlerContainer) AuthHandler() *handlers.AuthHandler {
	return c.authHandler
}
