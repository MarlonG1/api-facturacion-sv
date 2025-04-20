package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/invalidation_models"
)

type InvalidationReasonStrategy struct {
	Document *invalidation_models.InvalidationDocument
}

func (s *InvalidationReasonStrategy) Validate() *dte_errors.DTEError {
	reason := s.Document.Reason
	if reason == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Reason")
	}

	// Validar campos requeridos
	if reason.ResponsibleName == "" ||
		reason.ResponsibleDocType.GetValue() == "" ||
		reason.ResponsibleDocNum.GetValue() == "" ||
		reason.RequesterName == "" ||
		reason.RequesterDocType.GetValue() == "" ||
		reason.RequesterDocNum.GetValue() == "" {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Reason required fields")
	}

	// Validar enum de tipo anulación (1,2,3)
	if reason.Type.GetValue() < 1 || reason.Type.GetValue() > 3 {
		return dte_errors.NewDTEErrorSimple("InvalidEnum", "Annulment type")
	}

	// Validar tipos de documento
	if !contains(reason.ResponsibleDocType.GetValue()) ||
		!contains(reason.RequesterDocType.GetValue()) {
		return dte_errors.NewDTEErrorSimple("InvalidEnum", "Document type")
	}

	// Si es tipo 3, motivo es obligatorio
	if reason.Type.GetValue() == 3 && reason.Reason == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Reason for type 3")
	}

	// Si es tipo 2, codigoGeneracionR debe ser null
	if reason.Type.GetValue() == 2 {
		if s.Document.Document.ReplacementCode != nil {
			return dte_errors.NewDTEErrorSimple("InvalidField", "ReplacementCode must be null for type 2")
		}
	} else {
		// Para tipo 1 y 3, codigoGeneracionR es obligatorio
		if s.Document.Document.ReplacementCode == nil {
			return dte_errors.NewDTEErrorSimple("RequiredField", "ReplacementCode for types 1 and 3")
		}
	}

	return nil
}

func contains(item string) bool {
	for _, s := range constants.ValidReceiverDTETypes {
		if s == item {
			return true
		}
	}
	return false
}
