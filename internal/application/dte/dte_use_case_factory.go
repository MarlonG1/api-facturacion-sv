package dte

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation"
	domainPort "github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper"
)

// DTEUseCaseFactory facilita la creación de casos de uso para diferentes tipos de DTE
type DTEUseCaseFactory struct {
	authService       auth.AuthManager
	dteService        dte_documents.DTEManager
	transmitter       ports.BaseTransmitter
	mapperFactory     *mapper.MapperFactory
	operationsFactory *DTEOperations
}

// NewDTEUseCaseFactory crea una nueva instancia de DTEUseCaseFactory
func NewDTEUseCaseFactory(
	authService auth.AuthManager,
	dteService dte_documents.DTEManager,
	transmitter ports.BaseTransmitter,
) *DTEUseCaseFactory {
	return &DTEUseCaseFactory{
		authService:       authService,
		dteService:        dteService,
		transmitter:       transmitter,
		mapperFactory:     mapper.NewMapperFactory(),
		operationsFactory: NewDTEOperations(),
	}
}

// CreateInvoiceUseCase crea un caso de uso para facturas
func (f *DTEUseCaseFactory) CreateInvoiceUseCase(invoiceService domainPort.DTEService) *GenericDTEUseCase {
	return NewGenericDTEUseCase(
		f.authService,
		f.dteService,
		f.transmitter,
		invoiceService,
		f.mapperFactory.CreateInvoiceMapperAdapter(),
		f.mapperFactory.GetInvoiceResponseMapper(),
		f.operationsFactory.GetNoOperation(),
	)
}

// CreateCCFUseCase crea un caso de uso para CCF
func (f *DTEUseCaseFactory) CreateCCFUseCase(ccfService domainPort.DTEService) *GenericDTEUseCase {
	return NewGenericDTEUseCase(
		f.authService,
		f.dteService,
		f.transmitter,
		ccfService,
		f.mapperFactory.CreateCCFMapperAdapter(),
		f.mapperFactory.GetCCFResponseMapper(),
		f.operationsFactory.GetNoOperation(),
	)
}

// CreateCreditNoteUseCase crea un caso de uso para notas de crédito
func (f *DTEUseCaseFactory) CreateCreditNoteUseCase(creditNoteService domainPort.DTEService) *GenericDTEUseCase {
	return NewGenericDTEUseCase(
		f.authService,
		f.dteService,
		f.transmitter,
		creditNoteService,
		f.mapperFactory.CreateCreditNoteMapperAdapter(),
		f.mapperFactory.GetCreditNoteResponseMapper(),
		f.operationsFactory.GetCreditNoteOperations(f.dteService),
	)
}

// CreateRetentionUseCase crea un caso de uso para retenciones
func (f *DTEUseCaseFactory) CreateRetentionUseCase(retentionService domainPort.DTEService) *GenericDTEUseCase {
	return NewGenericDTEUseCase(
		f.authService,
		f.dteService,
		f.transmitter,
		retentionService,
		f.mapperFactory.CreateRetentionMapperAdapter(),
		f.mapperFactory.GetRetentionResponseMapper(),
		f.operationsFactory.GetNoOperation(),
	)
}

func (f *DTEUseCaseFactory) CreateInvalidationUseCase(
	invalidationManager invalidation.InvalidationManager,
) *InvalidationUseCase {
	return NewInvalidationUseCase(
		f.dteService,
		invalidationManager,
		f.authService,
		f.transmitter,
	)
}
