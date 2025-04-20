package document

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ContingencyType struct {
	Value int `json:"value"`
}

func NewContingencyType(value int) (*ContingencyType, error) {
	ct := &ContingencyType{Value: value}
	if ct.IsValid() {
		return ct, nil
	}
	return &ContingencyType{}, dte_errors.NewValidationError("InvalidLength", "contingency type", "1-5", fmt.Sprintf("%d", value))
}

func NewValidatedContingencyType(value int) *ContingencyType {
	return &ContingencyType{Value: value}
}

// IsValid valida que el valor de ContingencyType sea 1 a 5
func (ct *ContingencyType) IsValid() bool {
	return ct.Value >= 1 && ct.Value <= 5
}

func (ct *ContingencyType) Equals(other interfaces.ValueObject[int]) bool {
	return ct.GetValue() == other.GetValue()
}

func (ct *ContingencyType) GetValue() int {
	return ct.Value
}

func (ct *ContingencyType) ToString() string {
	return fmt.Sprintf("%d", ct.Value)
}
