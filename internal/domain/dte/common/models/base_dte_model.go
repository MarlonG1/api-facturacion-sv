package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator"
)

// DTEDocument es una estructura que representa un documento DTE, contiene Identification, Issuer, Receiver, Items, Summary, Extension y Appendix
type DTEDocument struct {
	Identification   interfaces.Identification    `json:"identification"`
	Issuer           interfaces.Issuer            `json:"issuer"`
	Receiver         interfaces.Receiver          `json:"receiver"`
	Items            []interfaces.Item            `json:"items"`
	Summary          interfaces.Summary           `json:"summary"`
	Extension        interfaces.Extension         `json:"extension,omitempty"`
	Appendix         []interfaces.Appendix        `json:"appendix,omitempty"`
	RelatedDocuments []interfaces.RelatedDocument `json:"related_documents,omitempty"`
	OtherDocuments   []interfaces.OtherDocuments  `json:"other_documents,omitempty"`
	ThirdPartySale   interfaces.ThirdPartySale    `json:"third_party_sale,omitempty"`
}

func (d *DTEDocument) GetIdentification() interfaces.Identification {
	return d.Identification
}
func (d *DTEDocument) GetIssuer() interfaces.Issuer {
	return d.Issuer
}
func (d *DTEDocument) GetReceiver() interfaces.Receiver {
	return d.Receiver
}
func (d *DTEDocument) GetItems() []interfaces.Item {
	return d.Items
}
func (d *DTEDocument) GetSummary() interfaces.Summary {
	return d.Summary
}
func (d *DTEDocument) GetAppendix() []interfaces.Appendix {
	return d.Appendix
}
func (d *DTEDocument) GetExtension() interfaces.Extension {
	return d.Extension
}
func (d *DTEDocument) GetRelatedDocuments() []interfaces.RelatedDocument {
	return d.RelatedDocuments
}
func (d *DTEDocument) GetOtherDocuments() []interfaces.OtherDocuments {
	return d.OtherDocuments
}
func (d *DTEDocument) GetThirdPartySale() interfaces.ThirdPartySale {
	return d.ThirdPartySale
}

// Validate Válida un documento DTE contra las reglas de validación
func (d *DTEDocument) Validate() error {
	return validator.ValidateDTEDocument(d)
}

// ValidateDTERules Válida las reglas de negocio de un documento DTE y retorna un error si no cumple con las reglas
func (d *DTEDocument) ValidateDTERules() *dte_errors.DTEError {
	rulesValidator := validator.NewDTERulesValidator(d)
	return rulesValidator.Validate()
}
