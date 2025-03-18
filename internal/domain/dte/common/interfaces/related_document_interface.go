package interfaces

import "time"

// RelatedDocument es una interfaz que define los metodos que deben ser implementados por los documentos relacionados
type RelatedDocument interface {
	GetDocumentType() string    // GetDocumentType retorna el tipo de documento relacionado
	GetGenerationType() int     // GetGenerationType retorna el tipo de generación del documento relacionado
	GetDocumentNumber() string  // GetDocumentNumber retorna el número de documento relacionado
	GetEmissionDate() time.Time // GetEmissionDate retorna la fecha de emisión del documento relacionado
}
