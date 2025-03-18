package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ExtensionStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas de extensión de un documento DTE
func (s *ExtensionStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetSummary() == nil {
		return nil
	}

	// Si el total de operaciones es mayor o igual a 1095, se requiere extensión
	if s.Document.GetSummary().GetTotalOperation() >= 1095.00 {
		if s.Document.GetExtension() == nil {
			return dte_errors.NewDTEErrorSimple("RequiredExtension", s.Document.GetSummary().GetTotalOperation())
		}
	}
	return nil
}
