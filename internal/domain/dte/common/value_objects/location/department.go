package location

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Department struct {
	Value string `json:"value"`
}

func NewDepartment(value string) (*Department, error) {
	dept := &Department{Value: value}
	if dept.IsValid() {
		return dept, nil
	}
	return &Department{}, dte_errors.NewValidationError("InvalidPattern", "Department", "01 a 14, debe ser de dos dígitos", value)
}

func NewValidatedDepartment(value string) *Department {
	return &Department{Value: value}
}

// IsValid válida que el valor de Department sea un número entre 01 y 14 (dos dígitos)
func (d *Department) IsValid() bool {
	pattern := `^0[1-9]|1[0-4]$`
	matched, _ := regexp.MatchString(pattern, d.Value)
	return matched
}

func (d *Department) Equals(other interfaces.ValueObject[string]) bool {
	return d.GetValue() == other.GetValue()
}

func (d *Department) GetValue() string {
	return d.Value
}

func (d *Department) ToString() string {
	return d.Value
}
