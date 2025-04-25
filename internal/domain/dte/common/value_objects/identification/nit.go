package identification

import (
	"regexp"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type NIT struct {
	Value string `json:"value"`
}

func NewNIT(value string) (*NIT, error) {
	//En caso de que el NIT tenga guiones, se eliminan para validar el patrón
	if strings.Contains(value, "-") {
		value = strings.ReplaceAll(value, "-", "")
	}
	value = strings.TrimSpace(value)

	nit := &NIT{Value: value}
	if nit.IsValid() {
		return nit, nil
	}
	return &NIT{}, dte_errors.NewValidationError("InvalidPattern", "nit", "12345678901234 o 123456789", value)
}

func NewValidatedNIT(value string) *NIT {
	return &NIT{Value: value}
}

// IsValid válido si el NIT tiene 14 o 9 dígitos
func (n *NIT) IsValid() bool {
	pattern := `^([0-9]{14}|[0-9]{9})$`
	matched, _ := regexp.MatchString(pattern, n.Value)
	return matched
}

func (n *NIT) Equals(other interfaces.ValueObject[string]) bool {
	return n.GetValue() == other.GetValue()
}

func (n *NIT) GetValue() string {
	return n.Value
}

func (n *NIT) ToString() string {
	return n.Value
}
