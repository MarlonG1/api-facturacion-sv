package bootstrap

import "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"

type HandlerContainer struct {
	useCases *UseCaseContainer

	authHandler    *handlers.AuthHandler
	invoiceHandler *handlers.InvoiceHandler
}

func NewHandlerContainer(useCases *UseCaseContainer) *HandlerContainer {
	return &HandlerContainer{
		useCases: useCases,
	}
}

func (c *HandlerContainer) Initialize() {
	c.authHandler = handlers.NewAuthHandler(c.useCases.AuthUseCase())
	c.invoiceHandler = handlers.NewInvoiceHandler(c.useCases.InvoiceUseCase())
}

func (c *HandlerContainer) AuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *HandlerContainer) InvoiceHandler() *handlers.InvoiceHandler {
	return c.invoiceHandler
}
