package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type DocumentTypeStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas de tipo de documento de un documento DTE
func (s *DocumentTypeStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetIdentification() == nil {
		return nil
	}

	docType := s.Document.GetIdentification().GetDTEType()

	switch docType {
	case constants.CCFElectronico:
		if s.Document.GetReceiver().GetNRC() == nil {
			return dte_errors.NewDTEErrorSimple("MissingNRC",
				s.Document.GetReceiver().GetName(),
				constants.CCFElectronico)
		}
	}

	return nil
}
