package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ModelTypeStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas de tipo de modelo de un documento DTE
func (s *ModelTypeStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetIdentification() == nil {
		return nil
	}

	// Si es una transmisión normal, el modelo debe ser de facturación previa
	if s.Document.GetIdentification().GetOperationType() == constants.TransmisionNormal &&
		s.Document.GetIdentification().GetModelType() != constants.ModeloFacturacionPrevio {
		return dte_errors.NewDTEErrorSimple("InvalidModelType",
			s.Document.GetIdentification().GetModelType())
	}

	// Si es una transmisión de contingencia, el modelo debe ser de invoice diferida
	if s.Document.GetIdentification().GetOperationType() == constants.TransmisionContingencia &&
		s.Document.GetIdentification().GetModelType() != constants.ModeloFacturacionDiferido {
		return dte_errors.NewDTEErrorSimple("InvalidModelType",
			s.Document.GetIdentification().GetModelType())
	}

	return nil
}
