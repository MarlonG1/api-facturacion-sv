package document

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type TransmissionType struct {
	Value int `json:"value"`
}

func NewTransmissionType(value int) (*TransmissionType, error) {
	tt := &TransmissionType{Value: value}
	if tt.IsValid() {
		return tt, nil
	}
	return &TransmissionType{}, dte_errors.NewValidationError("InvalidTransmissionType", fmt.Sprintf("%d", value))
}

func NewValidatedTransmissionType(value int) *TransmissionType {
	return &TransmissionType{Value: value}
}

// IsValid valida que el valor de TransmissionType sea 1 o 2 (1: Normal, 2: Contingencia)
func (tt *TransmissionType) IsValid() bool {
	return tt.Value == 1 || tt.Value == 2
}

func (tt *TransmissionType) Equals(other interfaces.ValueObject[int]) bool {
	return tt.GetValue() == other.GetValue()
}

func (tt *TransmissionType) GetValue() int {
	return tt.Value
}

func (tt *TransmissionType) ToString() string {
	return fmt.Sprintf("%d", tt.Value)
}
