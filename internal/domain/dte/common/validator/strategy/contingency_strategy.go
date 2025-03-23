package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ContingencyStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas de contingencia de un documento DTE
func (s *ContingencyStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetIdentification() == nil {
		return nil
	}

	if s.Document.GetIdentification().GetOperationType() == constants.TransmisionContingencia {
		// Si es transmisión por contingencia, debemos tener un tipo de contingencia válido (> 0)
		if s.Document.GetIdentification().GetContingencyType() == nil {
			return dte_errors.NewDTEErrorSimple("MissingContingencyType")
		}

		// Si es otro motivo (5), requiere razón
		if *s.Document.GetIdentification().GetContingencyType() == constants.OtroMotivo &&
			len(*s.Document.GetIdentification().GetContingencyReason()) == 0 {
			return dte_errors.NewDTEErrorSimple("MissingContingencyReason")
		}

		// Validar que sea uno de los tipos permitidos
		if !validateContingencyType(*s.Document.GetIdentification().GetContingencyType()) {
			return dte_errors.NewDTEErrorSimple("InvalidContingencyType",
				s.Document.GetIdentification().GetContingencyType())
		}
	}

	return nil
}

// validateContingencyType Válida que el tipo de contingencia sea uno de los permitidos
func validateContingencyType(contingencyType int) bool {
	for _, ct := range constants.AllowedContingencyTypes {
		if ct == contingencyType {
			return true
		}
	}
	return false
}
