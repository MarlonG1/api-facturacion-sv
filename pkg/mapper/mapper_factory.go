package mapper

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"time"
)

// MapperFactory es una fábrica para crear adaptadores de mappers
type MapperFactory struct{}

// NewMapperFactory crea una nueva instancia de MapperFactory
func NewMapperFactory() *MapperFactory {
	return &MapperFactory{}
}

// CreateInvoiceMapperAdapter crea un adaptador para el mapper de facturas
func (f *MapperFactory) CreateInvoiceMapperAdapter() DTEMapper {
	invoiceMapper := request_mapper.NewInvoiceMapper()

	return &MapperAdapter{
		MapFunc: func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
			invoiceReq, ok := req.(*structs.CreateInvoiceRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type, expected *structs.CreateInvoiceRequest")
			}
			return invoiceMapper.MapToInvoiceData(invoiceReq, issuer)
		},
	}
}

// CreateCCFMapperAdapter crea un adaptador para el mapper de CCF
func (f *MapperFactory) CreateCCFMapperAdapter() DTEMapper {
	ccfMapper := request_mapper.NewCCFMapper()

	return &MapperAdapter{
		MapFunc: func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
			ccfReq, ok := req.(*structs.CreateCreditFiscalRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type, expected *structs.CreateCreditFiscalRequest")
			}
			return ccfMapper.MapToCCFData(ccfReq, issuer)
		},
	}
}

// CreateCreditNoteMapperAdapter crea un adaptador para el mapper de Notas de Crédito
func (f *MapperFactory) CreateCreditNoteMapperAdapter() DTEMapper {
	creditNoteMapper := request_mapper.NewCreditNoteMapper()

	return &MapperAdapter{
		MapFunc: func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
			creditNoteReq, ok := req.(*structs.CreateCreditNoteRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type, expected *structs.CreateCreditNoteRequest")
			}
			return creditNoteMapper.MapToCreditNoteData(creditNoteReq, issuer)
		},
	}
}

// CreateRetentionMapperAdapter crea un adaptador para el mapper de Retenciones
func (f *MapperFactory) CreateRetentionMapperAdapter() DTEMapper {
	retentionMapper := request_mapper.NewRetentionMapper()

	return &MapperAdapter{
		MapFunc: func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
			retentionReq, ok := req.(*structs.CreateRetentionRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type, expected *structs.CreateRetentionRequest")
			}
			return retentionMapper.MapToRetentionData(retentionReq, issuer)
		},
	}
}

// CreateInvalidationMapperAdapter crea un adaptador para el mapper de Invalidaciones
func (f *MapperFactory) CreateInvalidationMapperAdapter() DTEMapper {
	invalidationMapper := request_mapper.NewInvalidationMapper()

	return &MapperAdapter{
		MapFunc: func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
			baseDte := params[0].(*dte.DTEDetails)
			emissionDate := params[1].(time.Time)

			invalidationReq, ok := req.(*structs.CreateInvalidationRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type, expected *structs.CreateInvalidationRequest")
			}
			return invalidationMapper.MapToInvalidationData(invalidationReq, issuer, baseDte, emissionDate)
		},
	}
}

// GetInvoiceResponseMapper devuelve la función de mapeo para respuestas de facturas
func (f *MapperFactory) GetInvoiceResponseMapper() ResponseMapperFunc {
	return func(domain interface{}) interface{} {
		return response_mapper.ToMHInvoice(domain)
	}
}

// GetCCFResponseMapper devuelve la función de mapeo para respuestas de CCF
func (f *MapperFactory) GetCCFResponseMapper() ResponseMapperFunc {
	return func(domain interface{}) interface{} {
		return response_mapper.ToMHCreditFiscalInvoice(domain)
	}
}

// GetCreditNoteResponseMapper devuelve la función de mapeo para respuestas de Notas de Crédito
func (f *MapperFactory) GetCreditNoteResponseMapper() ResponseMapperFunc {
	return func(domain interface{}) interface{} {
		return response_mapper.ToMHCreditNote(domain)
	}
}

// GetRetentionResponseMapper devuelve la función de mapeo para respuestas de Retenciones
func (f *MapperFactory) GetRetentionResponseMapper() ResponseMapperFunc {
	return func(domain interface{}) interface{} {
		return response_mapper.ToMHRetention(domain)
	}
}

// GetInvalidationResponseMapper devuelve la función de mapeo para respuestas de Invalidaciones
func (f *MapperFactory) GetInvalidationResponseMapper() ResponseMapperFunc {
	return func(domain interface{}) interface{} {
		return response_mapper.ToMHInvalidation(domain)
	}
}
