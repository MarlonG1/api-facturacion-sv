package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
)

// Appendix es una estructura que representa un anexo, contiene Field, Label y Value de un DTE
type Appendix struct {
	Field document.AppendixField `json:"field"`
	Label document.AppendixLabel `json:"label"`
	Value document.AppendixValue `json:"value"`
}

func (a *Appendix) GetField() string {
	return a.Field.GetValue()
}

func (a *Appendix) GetLabel() string {
	return a.Label.GetValue()
}

func (a *Appendix) GetValue() string {
	return a.Value.GetValue()
}
