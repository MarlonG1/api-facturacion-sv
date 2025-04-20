package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type CCFRelatedDocStrategy struct {
	Document *ccf_models.CreditFiscalDocument
}

// Validate - Valida los documentos relacionados de una Comprobante de Crédito Fiscal
func (s *CCFRelatedDocStrategy) Validate() *dte_errors.DTEError {
	if s.Document.GetRelatedDocuments() == nil || len(s.Document.GetRelatedDocuments()) == 0 {
		return nil
	}

	// No debe exceder el máximo de documentos relacionados
	if len(s.Document.GetRelatedDocuments()) > 50 {
		return dte_errors.NewDTEErrorSimple("ExceededRelatedDocsLimit",
			len(s.Document.GetRelatedDocuments()))
	}

	// Validar tipos de documentos relacionados permitidos para Comprobante de Crédito Fiscal
	for _, doc := range s.Document.GetRelatedDocuments() {
		if err := s.validateRelatedDocType(doc.GetDocumentType()); err != nil {
			return err
		}
	}

	// Validar consistencia de referencias en ítems
	for _, item := range s.Document.CreditItems {
		if item.GetRelatedDoc() == nil {
			logs.Error("Missing related document in item", map[string]interface{}{
				"itemNumber": item.GetNumber(),
			})
			return dte_errors.NewDTEErrorSimple("MissingItemRelatedDoc", item.GetNumber())
		}

		// Verificar que el documento relacionado del ítem exista en la lista de documentos relacionados
		found := false
		itemRelatedDoc := *item.GetRelatedDoc()

		for _, relDoc := range s.Document.GetRelatedDocuments() {
			if relDoc.GetDocumentNumber() == itemRelatedDoc {
				found = true
				break
			}
		}

		if !found {
			logs.Error("Item related document not found in document related docs", map[string]interface{}{
				"itemNumber": item.GetNumber(),
				"relatedDoc": itemRelatedDoc,
			})
			return dte_errors.NewDTEErrorSimple("InvalidItemRelatedDoc",
				item.GetNumber(),
				itemRelatedDoc)
		}
	}

	return nil
}

// validateRelatedDocType - Valida que el tipo de documento relacionado sea válido para Comprobante de Crédito Fiscal
func (s *CCFRelatedDocStrategy) validateRelatedDocType(docType string) *dte_errors.DTEError {

	if !constants.ValidCCFDTETypesRelateDoc[docType] {
		return dte_errors.NewDTEErrorSimple("InvalidRelatedDocDTEType", docType, constants.ShowValidRelatedDocTypes(constants.ValidCCFDTETypesRelateDoc))
	}

	return nil
}
