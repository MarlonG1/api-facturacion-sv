package identification

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type NRC struct {
	Value string `json:"value"`
}

func NewNRC(value string) (*NRC, error) {
	nrc := &NRC{Value: value}
	if nrc.IsValid() {
		return nrc, nil
	}
	return &NRC{}, dte_errors.NewValidationError("InvalidFormat", "nrc", "12345678", value)
}

func NewValidatedNRC(value string) *NRC {
	return &NRC{Value: value}
}

// IsValid válido si el NRC tiene entre 1 y 8 dígitos
func (n *NRC) IsValid() bool {
	pattern := `^[0-9]{1,8}$`
	matched, _ := regexp.MatchString(pattern, n.Value)
	return matched
}

func (n *NRC) Equals(other interfaces.ValueObject[string]) bool {
	return n.GetValue() == other.GetValue()
}

func (n *NRC) GetValue() string {
	return n.Value
}

func (n *NRC) ToString() string {
	return n.Value
}
