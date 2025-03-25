package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
)

type HandlerContainer struct {
	useCases *UseCaseContainer
	services *ServicesContainer

	authHandler        *handlers.AuthHandler
	invoiceHandler     *handlers.InvoiceHandler
	ccfHandler         *handlers.CCFHandler
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
	c.authHandler = handlers.NewAuthHandler(c.useCases.AuthUseCase())
	c.invoiceHandler = handlers.NewInvoiceHandler(c.useCases.InvoiceUseCase(), c.contingencyHandler)
	c.ccfHandler = handlers.NewCCFHandler(c.useCases.CCFUseCase(), c.contingencyHandler)
}

func (c *HandlerContainer) AuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *HandlerContainer) InvoiceHandler() *handlers.InvoiceHandler {
	return c.invoiceHandler
}

func (c *HandlerContainer) CCFHandler() *handlers.CCFHandler {
	return c.ccfHandler
}
