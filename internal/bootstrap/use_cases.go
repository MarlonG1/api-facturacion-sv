package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/application/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
)

type UseCaseContainer struct {
	services *ServicesContainer

	authUseCase         *auth.AuthUseCase
	invoiceUseCase      *dte.InvoiceUseCase
	ccfUseCase          *dte.CCFUseCase
	retentionUseCase    *dte.RetentionUseCase
	dteConsult          *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	creditNoteUseCase   *dte.CreditNoteUseCase
	baseTransmitter     ports.BaseTransmitter
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
	c.dteConsult = dte.NewDTEConsultUseCase(c.services.DTEManager())
	c.invalidationUseCase = dte.NewInvalidationUseCase(c.services.DTEManager(), c.services.InvalidationManager(), c.services.AuthManager(), c.baseTransmitter)
	c.retentionUseCase = dte.NewRetentionUseCase(c.services.AuthManager(), c.services.RetentionManager(), c.services.DTEManager(), c.baseTransmitter)
	c.creditNoteUseCase = dte.NewCreditNoteUseCase(c.services.AuthManager(), c.services.CreditNoteManager(), c.baseTransmitter, c.services.DTEManager())
}

func (c *UseCaseContainer) CreditNoteUseCase() *dte.CreditNoteUseCase {
	return c.creditNoteUseCase
}

func (c *UseCaseContainer) RetentionUseCase() *dte.RetentionUseCase {
	return c.retentionUseCase
}

func (c *UseCaseContainer) InvalidationUseCase() *dte.InvalidationUseCase {
	return c.invalidationUseCase
}

func (c *UseCaseContainer) DTEConsultUseCase() *dte.DTEConsultUseCase {
	return c.dteConsult
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
