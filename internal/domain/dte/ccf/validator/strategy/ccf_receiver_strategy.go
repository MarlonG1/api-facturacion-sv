package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type CCFReceiverStrategy struct {
	Document *ccf_models.CreditFiscalDocument
}

// Validate - Valida los campos específicos de un receptor de CCF
func (s *CCFReceiverStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetReceiver() == nil {
		return nil
	}

	// CCF requiere NRC obligatoriamente según schema
	if s.Document.GetReceiver().GetNRC() == nil {
		return dte_errors.NewDTEErrorSimple("MissingNRC",
			utils.PointerToString(s.Document.GetReceiver().GetName()),
			constants.CCFElectronico)
	}

	// CCF requiere código y descripción de actividad económica
	if s.Document.GetReceiver().GetActivityCode() == nil ||
		s.Document.GetReceiver().GetActivityDescription() == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField",
			"ActivityCode and Description")
	}

	return nil
}
