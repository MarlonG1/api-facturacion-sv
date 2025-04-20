package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// CreateDefaultReasonRequest crea una razón de invalidación predeterminada válida
func CreateDefaultReasonRequest() *structs.ReasonRequest {
	reason := "Factura con datos incorrectos"

	return &structs.ReasonRequest{
		Type:               1, // Reemplazo
		ResponsibleName:    "Juan Responsable",
		ResponsibleDocType: "13", // DUI
		ResponsibleNumDoc:  "12345678-9",
		RequestorName:      "Ana Solicitante",
		RequestorDocType:   "13", // DUI
		RequestorNumDoc:    "98765432-1",
		Reason:             &reason,
	}
}

// CreateDefaultInvalidationRequest crea una solicitud de invalidación predeterminada válida
func CreateDefaultInvalidationRequest() *structs.CreateInvalidationRequest {
	replacementCode := "FF54E9DB-79C3-42CE-B432-EC522C97EFB9"

	return &structs.CreateInvalidationRequest{
		GenerationCode:            "AD54E9BB-79A3-42AE-B432-EC522C97EFB7",
		Reason:                    CreateDefaultReasonRequest(),
		ReplacementGenerationCode: &replacementCode,
	}
}

// CreateInvalidationWithInvalidType crea una solicitud de invalidación con tipo inválido
func CreateInvalidationWithInvalidType() *structs.CreateInvalidationRequest {
	req := CreateDefaultInvalidationRequest()
	req.Reason.Type = 99 // Tipo inválido
	return req
}

// CreateInvalidationTypeWithoutReason crea una solicitud de invalidación de tipo 3 sin razón
func CreateInvalidationTypeWithoutReason() *structs.CreateInvalidationRequest {
	req := CreateDefaultInvalidationRequest()
	req.Reason.Type = 3     // Invalidación definitiva
	req.Reason.Reason = nil // Sin razón
	return req
}

// CreateInvalidationType2WithReplacementCode crea una solicitud de invalidación de tipo 2 con código de reemplazo
func CreateInvalidationType2WithReplacementCode() *structs.CreateInvalidationRequest {
	req := CreateDefaultInvalidationRequest()
	req.Reason.Type = 2 // Anulación
	// No debería tener código de reemplazo pero lo tiene
	replacementCode := "FF54E9DB-79C3-42CE-B432-EC522C97EFB9"
	req.ReplacementGenerationCode = &replacementCode
	return req
}
