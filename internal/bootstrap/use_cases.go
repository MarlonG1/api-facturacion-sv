package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/application/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
)

type UseCaseContainer struct {
	services *ServicesContainer

	authUseCase     *auth.AuthUseCase
	invoiceUseCase  *dte.InvoiceUseCase
	ccfUseCase      *dte.CCFUseCase
	baseTransmitter ports.BaseTransmitter
}

func NewUseCaseContainer(services *ServicesContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() {
	c.authUseCase = auth.NewAuthUseCase(c.services.AuthManager(), c.services.CryptManager())
	c.baseTransmitter = dte.NewBaseTransmitter(c.services.TransmitterManager(), c.services.SignerManager())
	c.invoiceUseCase = dte.NewInvoiceUseCase(c.services.AuthManager(), c.services.InvoiceService(), c.baseTransmitter, c.services.DTEManager())
	c.ccfUseCase = dte.NewCCFUseCase(c.services.AuthManager(), c.services.CCFService(), c.baseTransmitter, c.services.DTEManager())
}

func (c *UseCaseContainer) CCFUseCase() *dte.CCFUseCase {
	return c.ccfUseCase
}

func (c *UseCaseContainer) AuthUseCase() *auth.AuthUseCase {
	return c.authUseCase
}

func (c *UseCaseContainer) InvoiceUseCase() *dte.InvoiceUseCase {
	return c.invoiceUseCase
}
