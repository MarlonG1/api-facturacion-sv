package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
)

// UpdateContingencyIdentification actualiza la identificación de contingencia en el JSON del DTE.
func UpdateContingencyIdentification(document interface{}, contiType *int8, reason *string) (map[string]interface{}, error) {
	// 1. Convertir cualquier struct a map[string]interface{} mediante Marshal y Unmarshal
	var dteDoc map[string]interface{}
	jsonBytes, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DTE JSON:  %w", err)
	}

	// 1.1 Luego convertimos ese JSON a map[string]interface{}
	if err = json.Unmarshal(jsonBytes, &dteDoc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DTE JSON: %w", err)
	}

	// 2. Actualizar la identificación de contingencia en el JSON del DTE
	if identifi, exist := dteDoc["identificacion"]; exist {
		identification, ok := identifi.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to assert identification as map[string]interface{}")
		}
		identification["tipoModelo"] = constants.ModeloFacturacionDiferido
		identification["tipoOperacion"] = constants.TransmisionContingencia
		identification["tipoContingencia"] = *contiType
		identification["motivoContin"] = *reason
	}

	return dteDoc, nil
}

// SetReceptionStampIntoAppendix añade el sello de recepción al apéndice del documento.
func SetReceptionStampIntoAppendix(document string, receptionStamp *string) (string, error) {
	// 1. Determinar el tipo de DTE
	dteInfo, err := ExtractAuxiliarIdentificationFromStringJSON(document)
	if err != nil || dteInfo.Identification.DTEType == "" {
		return "", fmt.Errorf("failed to determine DTE type: %w", err)
	}

	// 2. Mapear el documento a un HashMap
	var dteDoc map[string]interface{}
	if err := json.Unmarshal([]byte(document), &dteDoc); err != nil {
		return "", fmt.Errorf("failed to unmarshal DTE JSON: %w", err)
	}

	// 3. Verificar si existe el campo de apéndices
	if appendix, exist := dteDoc["apendice"]; exist {
		if appendix == nil {
			// Si no existe, crear un nuevo apéndice
			dteDoc["apendice"] = []map[string]interface{}{
				{
					"Campo":    "Datos del documento",
					"Etiqueta": "Sello de recepción",
					"Valor":    *receptionStamp,
				},
			}
		} else {
			// Si existe, añadir el sello de recepción al apéndice
			appendices := appendix.([]interface{})
			newAppendices := append(appendices, map[string]interface{}{
				"Campo":    "Datos del documento",
				"Etiqueta": "Sello de recepción",
				"Valor":    *receptionStamp,
			})
			dteDoc["apendice"] = newAppendices
		}
	}

	// 4. Convertir el HashMap a JSON
	jsonData, err := json.Marshal(dteDoc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DTE JSON: %w", err)
	}

	return string(jsonData), nil
}
