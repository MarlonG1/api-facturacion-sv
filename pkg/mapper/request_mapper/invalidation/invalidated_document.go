package invalidation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

func MapInvalidatedDocument(baseDTE *dte.DTEDetails, request *structs.CreateInvalidationRequest, emissionDate time.Time) (*models.InvalidatedDocument, error) {
	if baseDTE == nil {
		return nil, shared_error.NewFormattedGeneralServiceError("InvalidationMapper", "MapToInvalidatedDocument", "InvalidBaseDTE")
	}

	// Deserializar JSON del DTE original
	var dteData map[string]interface{}
	if err := json.Unmarshal([]byte(baseDTE.JSONData), &dteData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DTE data: %w", err)
	}

	// Extraer datos del receptor del DTE original
	docType, numDoc, name, email, phone, err := extractReceptorData(dteData)
	if err != nil {
		return nil, err
	}

	// Extraer monto IVA del resumen
	montoIVA, err := extractIVAAmount(dteData)
	if err != nil {
		return nil, err
	}

	// Crear documento invalidado
	doc := &models.InvalidatedDocument{
		Type:           *document.NewValidatedDTEType(baseDTE.DTEType),
		GenerationCode: *identification.NewValidatedGenerationCode(baseDTE.ID),
		ControlNumber:  *identification.NewValidatedControlNumber(baseDTE.ControlNumber),
		ReceptionStamp: *baseDTE.ReceptionStamp,
		EmissionDate:   *temporal.NewValidatedEmissionDate(emissionDate),
		IVAAmount:      financial.NewValidatedAmount(montoIVA),
	}

	// Campos opcionales
	if phone != nil {
		doc.Phone = base.NewValidatedPhone(*phone)
	}

	if name != nil {
		doc.Name = name
	}

	if docType != nil && numDoc != nil {
		doc.DocumentType = document.NewValidatedDTEType(*docType)
		doc.DocumentNumber = identification.NewValidatedDocumentNumber(*numDoc)
	}

	if email != nil {
		doc.Email = base.NewValidatedEmail(*email)
	}

	// Código de generación de reemplazo si aplica
	if request.ReplacementGenerationCode != nil {
		doc.ReplacementCode = identification.NewValidatedGenerationCode(*request.ReplacementGenerationCode)
	}

	return doc, nil
}

func extractReceptorData(dteData map[string]interface{}) (*string, *string, *string, *string, *string, error) {
	var docType, numDoc, name, phone, email *string

	receptor, ok := dteData["receptor"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, nil, nil, nil
	}

	if nombreValue, ok := receptor["nombre"].(string); ok {
		name = &nombreValue
	} else {
		name = new(string)
	}

	if telefonoValue, ok := receptor["telefono"].(string); ok {
		phone = &telefonoValue
	} else {
		phone = new(string)
	}

	if correoValue, ok := receptor["correo"].(string); ok {
		email = &correoValue
	} else {
		email = new(string)
	}

	if tipoDocValue, ok := receptor["tipoDocumento"].(string); ok {
		docType = &tipoDocValue
	} else {
		if nitValue, ok := receptor["nit"].(string); ok {
			nit := constants.NIT
			return &nit, &nitValue, name, email, phone, nil
		}
		docType = new(string)
	}

	if numDocValue, ok := receptor["numDocumento"].(string); ok {
		numDoc = &numDocValue
	} else {
		numDoc = new(string)
	}

	return docType, numDoc, name, email, phone, nil
}

func extractIVAAmount(dteData map[string]interface{}) (float64, error) {
	resumen, ok := dteData["resumen"].(map[string]interface{})
	if !ok {
		return 0, shared_error.NewGeneralServiceError("InvalidationMapper", "extractIVAAmount", "invalid resumen data in DTE", nil)
	}

	ivaAmount, ok := resumen["totalIva"].(float64)
	if ok {
		return ivaAmount, nil
	}

	tributos, ok := resumen["tributos"].([]interface{})
	if !ok {
		return 0, nil
	}

	// Buscar IVA (código 20)
	for _, t := range tributos {
		tributo := t.(map[string]interface{})
		if tributo["codigoTributo"] == constants.TaxIVA {
			return tributo["valorTributo"].(float64), nil
		}
	}

	return 0, nil
}
