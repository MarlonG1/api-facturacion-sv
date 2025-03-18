package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
)

// RelatedDocument representa un documento relacionado a la invoice
type RelatedDocument struct {
	DocumentType   document.DTEType      `json:"documentType"`
	GenerationType document.ModelType    `json:"generationType"`
	DocumentNumber string                `json:"documentNumber"`
	EmissionDate   temporal.EmissionDate `json:"emissionDate"`
}

func (r *RelatedDocument) GetDocumentType() string {
	return r.DocumentType.GetValue()
}

func (r *RelatedDocument) GetGenerationType() int {
	return r.GenerationType.GetValue()
}

func (r *RelatedDocument) GetDocumentNumber() string {
	return r.DocumentNumber
}

func (r *RelatedDocument) GetEmissionDate() time.Time {
	return r.EmissionDate.GetValue()
}
