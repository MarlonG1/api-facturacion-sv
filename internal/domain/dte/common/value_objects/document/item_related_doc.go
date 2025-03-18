package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ItemRelatedDoc struct {
	value string
}

func NewItemRelatedDoc(value string) (*ItemRelatedDoc, error) {
	adc := &ItemRelatedDoc{value: value}
	if adc.IsValid() {
		return adc, nil
	}
	return &ItemRelatedDoc{}, dte_errors.NewValidationError("Invalid")
}

func (ird *ItemRelatedDoc) GetValue() string {
	return ird.value
}

func (ird *ItemRelatedDoc) Equals(other interfaces.ValueObject[string]) bool {
	return ird.value == other.GetValue()
}

func (ird *ItemRelatedDoc) ToString() string {
	return ird.value
}

func (ird *ItemRelatedDoc) IsValid() bool {
	return len(ird.value) >= 1 && len(ird.value) <= 36
}
