package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// UpdateContingencyIdentification actualiza la identificación de contingencia en el JSON del DTE.
func UpdateContingencyIdentification(identification *structs.DTEIdentification, contiType *int8, reason *string) {
	identification.TipoModelo = constants.ModeloFacturacionDiferido
	identification.TipoOperacion = constants.TransmisionContingencia
	tipoContingencia := int(*contiType)
	identification.TipoContingencia = &tipoContingencia
	identification.MotivoContin = reason
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
			appendices := appendix.([]map[string]interface{})
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
