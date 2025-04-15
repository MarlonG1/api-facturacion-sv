package strategy

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type RelatedDocsStrategy struct {
	Document interfaces.DTEDocument
}

// Validate valida la estrategia de documentos relacionados del DTE
func (s *RelatedDocsStrategy) Validate() *dte_errors.DTEError {
	// Si no hay documentos relacionados, es válido
	if s.Document.GetRelatedDocuments() == nil || len(s.Document.GetRelatedDocuments()) == 0 {
		return nil
	}

	// No debe exceder el máximo de documentos relacionados
	if len(s.Document.GetRelatedDocuments()) > 50 {
		return dte_errors.NewDTEErrorSimple("ExceededRelatedDocsLimit",
			len(s.Document.GetRelatedDocuments()))
	}

	firstDocType := s.Document.GetRelatedDocuments()[0].GetDocumentType()
	for i, doc := range s.Document.GetRelatedDocuments() {
		if doc.GetDocumentType() != firstDocType {
			logs.Error("Mixed document types not allowed", map[string]interface{}{
				"expectedType": firstDocType,
				"foundType":    doc.GetDocumentType(),
				"index":        i,
			})
			return dte_errors.NewDTEErrorSimple("MixedDocumentTypesNotAllowed")
		}
	}

	// Validar cada documento relacionado
	for _, doc := range s.Document.GetRelatedDocuments() {
		if err := s.validateRelatedDoc(doc); err != nil {
			return err
		}
	}

	return nil
}

// validateRelatedDoc valida un documento relacionado
func (s *RelatedDocsStrategy) validateRelatedDoc(doc interfaces.RelatedDocument) *dte_errors.DTEError {

	// Validar que la fecha no sea futura
	if doc.GetEmissionDate().After(utils.TimeNow()) {
		return dte_errors.NewDTEErrorSimple("InvalidRelatedDocDate",
			doc.GetEmissionDate().Format("2006-01-02"))
	}

	// Validar formato de número de documento
	err := validateElectronicDocNumber(doc.GetDocumentNumber(), doc.GetGenerationType())
	if err != nil {
		return err
	}

	return nil
}

// validateElectronicDocNumber valida el número de documento electrónico
func validateElectronicDocNumber(number string, generationType int) *dte_errors.DTEError {

	if generationType == constants.TransmisionContingencia {
		if !isValidUUID(number) {
			return dte_errors.NewDTEErrorSimple("InvalidRelatedDocNumberContingency", number)
		}
	} else {
		if len(number) < 0 || len(number) > 20 {
			return dte_errors.NewDTEErrorSimple("InvalidRelatedDocNumberNormal", number)
		}
	}

	return nil
}

// isValidUUID valida que el número de documento relacionado sea un UUID válido
func isValidUUID(uuid string) bool {
	var uuidRegex = regexp.MustCompile(`^[A-F0-9]{8}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{12}$`)
	return uuidRegex.MatchString(uuid)
}
