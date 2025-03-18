package strategy

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

const (
	DUIPattern = `^[0-9]{8}-[0-9]{1}$`
	NITPattern = `^([0-9]{14}|[0-9]{9})$`
)

var (
	nitRegex = regexp.MustCompile(NITPattern)
	duiRegex = regexp.MustCompile(DUIPattern)
)

type BasicRulesStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas básicas de un documento DTE
func (s *BasicRulesStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil {
		return dte_errors.NewDTEErrorSimple("RequiredFieldMissing", "Document", "nil")
	}

	// Validación de identificación primero
	if s.Document.GetIdentification() == nil {
		return dte_errors.NewDTEErrorSimple("RequiredFieldMissing", "Identification", "nil")
	}

	docType := s.Document.GetIdentification().GetDTEType()

	// Validaciones subsecuentes que usan docType
	if s.Document.GetIssuer() == nil {
		return dte_errors.NewDTEErrorSimple("RequiredFieldMissing", "Issuer", docType)
	}

	if s.Document.GetItems() == nil || len(s.Document.GetItems()) == 0 {
		return dte_errors.NewDTEErrorSimple("RequiredFieldMissing", "Items", docType)
	}

	// Protección contra NPE al acceder al receiver
	if requiresReceiver(docType) {
		receiver := s.Document.GetReceiver()
		if receiver == nil {
			return dte_errors.NewDTEErrorSimple("RequiredFieldMissing", "Receiver", docType)
		}
	}

	err := getDocumentNumberError(s.Document.GetReceiver().GetDocumentNumber(), s.Document.GetReceiver().GetDocumentType())
	if err != nil {
		return err
	}

	return nil
}

// requiresReceiver Verifica si el tipo de documento requiere receptor
func requiresReceiver(docType string) bool {
	switch docType {
	case constants.FacturaElectronica,
		constants.CCFElectronico,
		constants.NotaRemisionElectronica,
		constants.NotaCreditoElectronica,
		constants.NotaDebitoElectronica,
		constants.FacturaExportacionElectronica:
		return true
	default:
		return false
	}
}

func getDocumentNumberError(documentNumber, documentType *string) *dte_errors.DTEError {
	switch *documentType {
	case constants.NIT:
		if !nitRegex.MatchString(*documentNumber) {
			return dte_errors.NewDTEErrorSimple("InvalidNITFormat", documentNumber)
		}
	case constants.DUI:
		if !duiRegex.MatchString(*documentNumber) {
			return dte_errors.NewDTEErrorSimple("InvalidDUIFormat", documentNumber)
		}
	}

	return nil
}
