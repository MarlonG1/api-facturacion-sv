package identification

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"strings"
)

// DocumentNumber representa el número de documento de identificación, un atributo opcional de un receptor
type DocumentNumber struct {
	Value string
}

func NewDocumentNumber(value string) (*DocumentNumber, error) {
	//En caso de que el numero de doc tenga guiones, se eliminan para validar el patrón
	if strings.Contains(value, "-") {
		value = strings.ReplaceAll(value, "-", "")
	}

	documentNumber := &DocumentNumber{Value: value}
	if documentNumber.IsValid() {
		return documentNumber, nil
	}
	return &DocumentNumber{}, dte_errors.NewValidationError("InvalidDocumentNumber", value)
}

func NewValidatedDocumentNumber(value string) *DocumentNumber {
	return &DocumentNumber{Value: value}
}

func (dn *DocumentNumber) IsValid() bool {
	return len(dn.Value) >= 3 && len(dn.Value) <= 20
}

func (dn *DocumentNumber) GetValue() string {
	return dn.Value
}

func (dn *DocumentNumber) Equals(other interfaces.ValueObject[string]) bool {
	return dn.Value == other.GetValue()
}

func (dn *DocumentNumber) ToString() string {
	return dn.Value
}
