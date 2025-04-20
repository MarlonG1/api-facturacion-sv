package mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
)

// DTEMapper es una interfaz genérica para todos los mappers de DTE
type DTEMapper interface {
	// MapToDomainModel convierte un request y un issuer a un modelo de dominio
	MapToDomainModel(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error)
}

// MapperAdapter adapta los mappers existentes a la interfaz DTEMapper
type MapperAdapter struct {
	// La función que implementa el mapeo real
	MapFunc func(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error)
}

// MapToDomainModel implementa la interfaz DTEMapper
func (a *MapperAdapter) MapToDomainModel(req interface{}, issuer *dte.IssuerDTE, params ...interface{}) (interface{}, error) {
	return a.MapFunc(req, issuer, params)
}

// ResponseMapperFunc define el tipo para funciones de mapeo de respuesta
type ResponseMapperFunc func(domain interface{}) interface{}
