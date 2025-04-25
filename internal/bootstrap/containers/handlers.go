package containers

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
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
	c.dteHandler = handlers.NewDTEHandler(c.useCases.DTEConsultUseCase(), c.useCases.InvalidationUseCase(),
		c.initializeGenericCreatorHandler(c.contingencyHandler),
	)
}

func (c *HandlerContainer) initializeGenericCreatorHandler(contingencyHandler *helpers.ContingencyHandler) *handlers.GenericCreatorDTEHandler {
	// Crear el handler genérico
	genericHandler := handlers.NewGenericDTEHandler(contingencyHandler)

	// Registrar los tipos de documentos
	genericHandler.RegisterDocument("/dte/invoices", helpers.DocumentConfig{
		UseCase:         c.useCases.InvoiceUseCase(),
		RequestType:     &structs.CreateInvoiceRequest{},
		DocumentType:    constants.FacturaElectronica,
		UsesContingency: true,
	})

	genericHandler.RegisterDocument("/dte/ccf", helpers.DocumentConfig{
		UseCase:         c.useCases.CCFUseCase(),
		RequestType:     &structs.CreateCreditFiscalRequest{},
		DocumentType:    constants.CCFElectronico,
		UsesContingency: true,
	})

	genericHandler.RegisterDocument("/dte/creditnote", helpers.DocumentConfig{
		UseCase:         c.useCases.CreditNoteUseCase(),
		RequestType:     &structs.CreateCreditNoteRequest{},
		DocumentType:    constants.NotaCreditoElectronica,
		UsesContingency: false, // TODO: activar cuando Hacienda resuelva el problema
	})

	genericHandler.RegisterDocument("/dte/retention", helpers.DocumentConfig{
		UseCase:         c.useCases.RetentionUseCase(),
		RequestType:     &structs.CreateRetentionRequest{},
		DocumentType:    constants.ComprobanteRetencionElectronico,
		UsesContingency: false,
	})

	return genericHandler
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
