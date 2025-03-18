package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type AppendixField struct {
	Value string `json:"value"`
}

func NewAppendixField(value string) (*AppendixField, error) {
	field := &AppendixField{Value: value}
	if field.IsValid() {
		return field, nil
	}
	return &AppendixField{}, dte_errors.NewValidationError("InvalidLength", "AppendixField", "2 a 25", value)
}

func NewValidatedAppendixField(value string) *AppendixField {
	return &AppendixField{Value: value}
}

func (f *AppendixField) IsValid() bool {
	return len(f.Value) >= 2 && len(f.Value) <= 25
}

func (f *AppendixField) Equals(other interfaces.ValueObject[string]) bool {
	return f.Value == other.GetValue()
}

func (f *AppendixField) GetValue() string {
	return f.Value
}

func (f *AppendixField) ToString() string {
	return f.Value
}
