package financial

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type PaymentType struct {
	Value string `json:"value"`
}

func NewPaymentType(value string) (*PaymentType, error) {
	pt := &PaymentType{Value: value}
	if pt.IsValid() {
		return pt, nil
	}
	return &PaymentType{}, dte_errors.NewValidationError("InvalidLength", "PaymentType", "01 a 14, 99", value)
}

func NewValidatedPaymentType(value string) *PaymentType {
	return &PaymentType{Value: value}
}

// IsValid válida que el valor de PaymentType sea 01 a 14, 99 (dos dígitos)
func (pt *PaymentType) IsValid() bool {
	pattern := `^(0[1-9]||1[0-4]||99)$`
	matched, _ := regexp.MatchString(pattern, pt.Value)
	return matched
}

func (pt *PaymentType) Equals(other interfaces.ValueObject[string]) bool {
	return pt.GetValue() == other.GetValue()
}

func (pt *PaymentType) GetValue() string {
	return pt.Value
}

func (pt *PaymentType) ToString() string {
	return pt.Value
}
