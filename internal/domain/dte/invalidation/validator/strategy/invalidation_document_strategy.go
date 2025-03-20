package strategy

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type InvalidationDocumentStrategy struct {
	Document *models.InvalidationDocument
}

func (s *InvalidationDocumentStrategy) Validate() *dte_errors.DTEError {
	doc := s.Document.Document
	if doc == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Document to invalidate")
	}

	// 1. Validar campos obligatorios según schema
	if doc.Type.GetValue() == "" || doc.GenerationCode.GetValue() == "" || doc.ReceptionStamp == "" ||
		doc.ControlNumber.GetValue() == "" || doc.EmissionDate.GetValue().IsZero() {
		logs.Info("DEBUG", map[string]interface{}{
			"docType":        doc.Type.GetValue(),
			"generationCode": doc.GenerationCode.GetValue(),
			"receptionStamp": doc.ReceptionStamp,
			"controlNumber":  doc.ControlNumber.GetValue(),
			"emissionDate":   doc.EmissionDate.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("RequiredField", "Document")
	}

	// 2. Validar ReceptionStamp patrón según schema
	if matched, _ := regexp.MatchString("^[A-Z0-9]{40}$", doc.ReceptionStamp); !matched {
		return dte_errors.NewDTEErrorSimple("InvalidPattern", "Reception stamp", "40 caracteres alfanumericos", doc.ReceptionStamp)
	}

	_, err := document.NewDTEType(doc.Type.GetValue())
	if err != nil {
		return dte_errors.NewDTEErrorSimple("InvalidDTETypeForInvalidation", doc.Type.GetValue())
	}

	return nil
}
