package document

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type OperationType struct {
	Value int `json:"value"`
}

func NewOperationType(value int) (*OperationType, error) {
	opType := &OperationType{Value: value}
	if opType.IsValid() {
		return opType, nil
	}
	return &OperationType{}, dte_errors.NewValidationError("InvalidLength", "OperationType", "1 o 2", fmt.Sprintf("%d", value))
}

func NewValidatedOperationType(value int) *OperationType {
	return &OperationType{Value: value}
}

// IsValid válida que el valor de OperationType sea 1 o 2
func (ot *OperationType) IsValid() bool {
	return ot.Value == 1 || ot.Value == 2
}

func (ot *OperationType) Equals(other interfaces.ValueObject[int]) bool {
	return ot.GetValue() == other.GetValue()
}

func (ot *OperationType) GetValue() int {
	return ot.Value
}

func (ot *OperationType) ToString() string {
	return fmt.Sprintf("%d", ot.Value)
}
