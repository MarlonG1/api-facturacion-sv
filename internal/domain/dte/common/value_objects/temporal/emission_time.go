package temporal

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type EmissionTime struct {
	Value time.Time `json:"value"`
}

func NewEmissionTime(value time.Time) (*EmissionTime, error) {
	timeValue := &EmissionTime{Value: value}
	if timeValue.IsValid() {
		return timeValue, nil
	}
	return &EmissionTime{}, dte_errors.NewValidationError("InvalidEmissionTime", value.String())
}

func NewValidatedEmissionTime(value time.Time) *EmissionTime {
	return &EmissionTime{Value: value}
}

// IsValid válido si la fecha no es cero
func (et *EmissionTime) IsValid() bool { return !et.Value.IsZero() }

func (et *EmissionTime) Equals(other interfaces.ValueObject[time.Time]) bool {
	return et.GetValue().Equal(other.GetValue())
}

func (et *EmissionTime) GetValue() time.Time {
	return et.Value
}

func (et *EmissionTime) ToString() string {
	return et.Value.Format("15:04:05")
}
