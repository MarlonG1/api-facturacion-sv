package document

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Version struct {
	Value int `json:"value"`
}

func NewVersion(value int) (*Version, error) {
	version := &Version{Value: value}
	if version.IsValid() {
		return version, nil
	}
	return &Version{}, dte_errors.NewValidationError("InvalidLength", "Version", "1, 2 o 3", fmt.Sprintf("%d", value))
}

func NewValidatedVersion(value int) *Version {
	return &Version{Value: value}
}

// IsValid valida que el valor de Version sea 1, 2 o 3
func (v *Version) IsValid() bool {
	return v.Value > 0 && v.Value <= 3
}

func (v *Version) Equals(other interfaces.ValueObject[int]) bool {
	return v.GetValue() == other.GetValue()
}

func (v *Version) GetValue() int {
	return v.Value
}

func (v *Version) ToString() string {
	return fmt.Sprintf("%d", v.Value)
}
