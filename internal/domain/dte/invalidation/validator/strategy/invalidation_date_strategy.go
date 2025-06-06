package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/invalidation_models"
)

type InvalidationDateStrategy struct {
	Document *invalidation_models.InvalidationDocument
}

func (s *InvalidationDateStrategy) Validate() *dte_errors.DTEError {
	if s.Document.Identification == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Identification")
	}

	// Validar plazos según tipo de documento
	docType := s.Document.Document.Type.GetValue()
	emissionDate := s.Document.Document.EmissionDate.GetValue()
	annulmentDate := s.Document.Identification.EmissionDate.GetValue()

	switch docType {
	case "01", "11": // Factura y FEXE: 3 meses
		if annulmentDate.Sub(emissionDate).Hours() > 24*90 {
			return dte_errors.NewDTEErrorSimple("InvalidDateForFEFX")
		}
	default: // Resto: 1 día
		if annulmentDate.Sub(emissionDate).Hours() > 24 {
			return dte_errors.NewDTEErrorSimple("InvalidDateForAllDTE")
		}
	}
	return nil
}
