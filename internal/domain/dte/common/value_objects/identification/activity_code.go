package identification

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ActivityCode struct {
	Value string `json:"value"`
}

func NewActivityCode(value string) (*ActivityCode, error) {
	code := &ActivityCode{Value: value}
	if code.IsValid() {
		return code, nil
	}
	return &ActivityCode{}, dte_errors.NewValidationError("InvalidPattern", "ActivityCode", "123456 o 12345678", value)
}

func NewValidatedActivityCode(value string) *ActivityCode {
	return &ActivityCode{Value: value}
}

// IsValid válida que el código de actividad económica tenga entre 2 y 6 dígitos
func (ac *ActivityCode) IsValid() bool {
	pattern := `^[0-9]{2,6}$`
	matched, _ := regexp.MatchString(pattern, ac.Value)
	return matched
}

func (ac *ActivityCode) Equals(other interfaces.ValueObject[string]) bool {
	return ac.GetValue() == other.GetValue()
}

func (ac *ActivityCode) GetValue() string {
	return ac.Value
}

func (ac *ActivityCode) ToString() string {
	return ac.Value
}
