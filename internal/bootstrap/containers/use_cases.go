package containers

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/application/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
)

type UseCaseContainer struct {
	services *ServicesContainer

	// Caso de uso especiales
	dteConsult          *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	authUseCase         *auth.AuthUseCase
	baseTransmitter     ports.BaseTransmitter
	dteUseCaseFactory   *dte.DTEUseCaseFactory

	// Casos de uso genéricos creacional
	invoiceUseCase    *dte.GenericDTEUseCase
	ccfUseCase        *dte.GenericDTEUseCase
	retentionUseCase  *dte.GenericDTEUseCase
	creditNoteUseCase *dte.GenericDTEUseCase
}

func NewUseCaseContainer(services *ServicesContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() {
	c.authUseCase = auth.NewAuthUseCase(c.services.AuthManager(), c.services.CryptManager())
	c.baseTransmitter = dte.NewBaseTransmitter(c.services.TransmitterManager(), c.services.SignerManager())
	c.dteConsult = dte.NewDTEConsultUseCase(c.services.DTEManager())

	// Inicializar factory de casos de uso
	c.dteUseCaseFactory = dte.NewDTEUseCaseFactory(
		c.services.AuthManager(),
		c.services.DTEManager(),
		c.baseTransmitter)

	c.invoiceUseCase = c.dteUseCaseFactory.CreateInvoiceUseCase(c.services.InvoiceService())
	c.ccfUseCase = c.dteUseCaseFactory.CreateCCFUseCase(c.services.CCFService())
	c.retentionUseCase = c.dteUseCaseFactory.CreateRetentionUseCase(c.services.RetentionManager())
	c.creditNoteUseCase = c.dteUseCaseFactory.CreateCreditNoteUseCase(c.services.CreditNoteManager())

	// Crear el caso de uso específico para invalidación
	c.invalidationUseCase = c.dteUseCaseFactory.CreateInvalidationUseCase(c.services.InvalidationManager())

}

func (c *UseCaseContainer) DTEConsultUseCase() *dte.DTEConsultUseCase {
	return c.dteConsult
}

func (c *UseCaseContainer) InvoiceUseCase() *dte.GenericDTEUseCase {
	return c.invoiceUseCase
}

func (c *UseCaseContainer) CCFUseCase() *dte.GenericDTEUseCase {
	return c.ccfUseCase
}

func (c *UseCaseContainer) RetentionUseCase() *dte.GenericDTEUseCase {
	return c.retentionUseCase
}

func (c *UseCaseContainer) CreditNoteUseCase() *dte.GenericDTEUseCase {
	return c.creditNoteUseCase
}

func (c *UseCaseContainer) InvalidationUseCase() *dte.InvalidationUseCase {
	return c.invalidationUseCase
}

func (c *UseCaseContainer) AuthUseCase() *auth.AuthUseCase {
	return c.authUseCase
}
