// strategy/other_documents_strategy.go

package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type OtherDocumentsStrategy struct {
	Document interfaces.DTEDocument
}

// Validate valida los documentos adicionales del DTE
func (s *OtherDocumentsStrategy) Validate() *dte_errors.DTEError {
	docs := s.Document.GetOtherDocuments()
	if docs == nil {
		return nil // Es opcional
	}

	// Validar cantidad de documentos
	if len(docs) == 0 || len(docs) > 10 {
		return dte_errors.NewDTEErrorSimple("InvalidOtherDocsCount", len(docs))
	}

	// Validar cada documento
	for _, doc := range docs {
		if err := s.validateDocument(doc); err != nil {
			return err
		}
	}

	return nil
}

// Validar documento
func (s *OtherDocumentsStrategy) validateDocument(doc interfaces.OtherDocuments) *dte_errors.DTEError {
	// Validar código
	code := doc.GetAssociatedDocument()
	if code < 1 || code > 4 {
		return dte_errors.NewDTEErrorSimple("InvalidAssociatedDocumentCode", code)
	}

	if code == 3 {
		return s.validateMedicalDocument(doc)
	}

	return s.validateRegularDocument(doc)
}

// Validar documento médico
func (s *OtherDocumentsStrategy) validateMedicalDocument(doc interfaces.OtherDocuments) *dte_errors.DTEError {
	if doc.GetDescription() != "" || doc.GetDetail() != "" {
		return dte_errors.NewDTEErrorSimple("InvalidMedicalDocFields")
	}

	if doc.GetDoctor() == nil {
		return dte_errors.NewDTEErrorSimple("RequiredField", "OtherDocuments->Doctor")
	}

	return s.validateDoctor(doc.GetDoctor())
}

// Validar documento regular (no médico)
func (s *OtherDocumentsStrategy) validateRegularDocument(doc interfaces.OtherDocuments) *dte_errors.DTEError {
	if doc.GetDoctor().GetName() != "" || doc.GetDoctor().GetServiceType() != 0 || doc.GetDoctor().GetNIT() != "" || doc.GetDoctor().GetIdentification() != "" {
		return dte_errors.NewDTEErrorSimple("InvalidField", "Doctor must be null")
	}

	if doc.GetAssociatedDocument() != 3 {
		if doc.GetDescription() == "" {
			return dte_errors.NewDTEErrorSimple("RequiredField", "OtherDocuments->Description")
		}

		if doc.GetDetail() == "" {
			return dte_errors.NewDTEErrorSimple("RequiredField", "OtherDocuments->Detail")
		}
	}

	if len(doc.GetDescription()) > 100 {
		return dte_errors.NewDTEErrorSimple("InvalidLength", "Description", "1-100", doc.GetDescription())
	}

	if len(doc.GetDetail()) > 300 {
		return dte_errors.NewDTEErrorSimple("InvalidLength", "Detail", "1-300", doc.GetDetail())
	}

	return nil
}

// validateDoctor valida los datos del doctor en un documento médico
func (s *OtherDocumentsStrategy) validateDoctor(doctor interfaces.DoctorInfo) *dte_errors.DTEError {
	if len(doctor.GetName()) == 0 || len(doctor.GetName()) > 100 {
		return dte_errors.NewDTEErrorSimple("InvalidLength", "DoctorName", "1-100", doctor.GetName())
	}

	serviceType := doctor.GetServiceType()
	if serviceType < 1 || serviceType > 6 {
		return dte_errors.NewDTEErrorSimple("InvalidServiceType", serviceType)
	}

	// Validar NIT y documento de identificación mutuamente excluyentes
	hasNIT := doctor.GetNIT() != ""
	hasID := doctor.GetIdentification() != ""

	if !hasNIT && !hasID {
		return dte_errors.NewDTEErrorSimple("RequiredField", "Doctor NIT or Identification")
	}

	if hasNIT && hasID {
		return dte_errors.NewDTEErrorSimple("MutuallyExclusiveFields", "NIT", "Identification")
	}

	return nil
}
