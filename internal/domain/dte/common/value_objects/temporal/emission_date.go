package temporal

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type EmissionDate struct {
	Value time.Time `json:"value"`
}

func NewEmissionDate(value time.Time) (*EmissionDate, error) {
	date := &EmissionDate{Value: value}
	if date.IsValid() {
		return date, nil
	}
	return &EmissionDate{}, dte_errors.NewValidationError("InvalidDateTime", value.String())
}

func NewValidatedEmissionDate(value time.Time) *EmissionDate {
	return &EmissionDate{Value: value}
}

// IsValid válido si la fecha no es cero y es menor a la fecha actual
func (ed *EmissionDate) IsValid() bool {
	return !ed.Value.IsZero() && ed.Value.Before(utils.TimeNow().Add(time.Hour*24))
}

func (ed *EmissionDate) Equals(other interfaces.ValueObject[time.Time]) bool {
	return ed.GetValue().Equal(other.GetValue())
}

func (ed *EmissionDate) GetValue() time.Time {
	return ed.Value
}

func (ed *EmissionDate) ToString() string {
	return ed.Value.Format("2006-01-02")
}
