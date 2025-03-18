package base

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Email struct {
	Value string `json:"value"`
}

func NewEmail(value string) (*Email, error) {
	email := &Email{Value: value}
	if email.IsValid() {
		return email, nil
	}
	return &Email{}, dte_errors.NewValidationError("InvalidPattern", "Email", "example@example.exam", value)
}

func NewValidatedEmail(value string) *Email {
	return &Email{Value: value}
}

// IsValid válido si el email cumple con el patrón de email
func (e *Email) IsValid() bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, e.Value)
	return matched && len(e.Value) >= 3 && len(e.Value) <= 100
}

func (e *Email) Equals(other interfaces.ValueObject[string]) bool {
	return e.GetValue() == other.GetValue()
}

func (e *Email) GetValue() string {
	return e.Value
}

func (e *Email) ToString() string {
	return e.Value
}
