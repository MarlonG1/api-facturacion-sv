package interfaces

import "time"

// RelatedDocumentGetter es una interfaz que define los métodos getter que deben ser implementados por los documentos relacionados
type RelatedDocumentGetter interface {
	GetDocumentType() string    // GetDocumentType retorna el tipo de documento relacionado
	GetGenerationType() int     // GetGenerationType retorna el tipo de generación del documento relacionado
	GetDocumentNumber() string  // GetDocumentNumber retorna el número de documento relacionado
	GetEmissionDate() time.Time // GetEmissionDate retorna la fecha de emisión del documento relacionado
}

// RelatedDocumentSetter es una interfaz que define los métodos setter que deben ser implementados por los documentos relacionados
type RelatedDocumentSetter interface {
	SetDocumentType(documentType string) error     // SetDocumentType establece el tipo de documento relacionado
	SetGenerationType(generationType int) error    // SetGenerationType establece el tipo de generación del documento relacionado
	SetDocumentNumber(documentNumber string) error // SetDocumentNumber establece el número de documento relacionado
	SetEmissionDate(emissionDate time.Time) error  // SetEmissionDate establece la fecha de emisión del documento relacionado
}

// RelatedDocument es una interfaz que combina los getters y setters de RelatedDocument
type RelatedDocument interface {
	RelatedDocumentGetter
	RelatedDocumentSetter
}
