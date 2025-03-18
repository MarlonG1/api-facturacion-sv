package strategy

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type TemporalValidationStrategy struct {
	Document interfaces.DTEDocument
}

// Validate valida la fecha y hora de emisión del documento
func (s *TemporalValidationStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetIdentification() == nil {
		return nil
	}

	emissionDate := s.Document.GetIdentification().GetEmissionDate()
	emissionTime := s.Document.GetIdentification().GetEmissionTime()
	now := utils.TimeNow()

	// Validar fecha futura
	if emissionDate.After(now) {
		return dte_errors.NewDTEErrorSimple("InvalidDateTime",
			emissionDate.Format("2006-01-02"))
	}

	// Validar hora futura en el mismo día
	if emissionDate.Equal(now.Truncate(24*time.Hour)) &&
		emissionTime.After(now) {
		return dte_errors.NewDTEErrorSimple("InvalidEmissionTime",
			emissionTime.Format("15:04:05"))
	}

	return nil
}
