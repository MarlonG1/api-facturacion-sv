package helpers

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/gorilla/mux"
	"net/http"
)

// DocumentConfig contiene la configuración para manejar un tipo específico de documento
type DocumentConfig struct {
	UseCase         *dte.GenericDTEUseCase
	RequestType     interface{}
	DocumentType    string
	UsesContingency bool
}

func GetRequestVar(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
