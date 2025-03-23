package processors

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"strconv"
	"strings"
)

func GetDocumentRequestData(document interface{}) (int, string, string, int, error) {
	var docMap map[string]interface{}

	jsonData, err := json.Marshal(document)
	if err != nil {
		return 0, "", "", 0, err
	}

	if err := json.Unmarshal(jsonData, &docMap); err != nil {
		return 0, "", "", 0, err
	}

	// Buscar el objeto identificacion
	identification, ok := docMap["identificacion"].(map[string]interface{})
	if !ok {
		logs.Error("Missing or invalid identificacion object")
		return 0, "", "", 0, fmt.Errorf("missing or invalid identificacion object")
	}

	// Extraer version
	version, ok := identification["version"].(float64)
	if !ok {
		logs.Error("Missing or invalid version field")
		return 0, "", "", 0, fmt.Errorf("missing or invalid version field")
	}

	dteType, ok := identification["tipoDte"].(string)
	if !ok {
		documento := docMap["documento"].(map[string]interface{})
		dteType, ok = documento["tipoDte"].(string)
		if !ok {
			logs.Error("Missing or invalid tipoDte field")
			return 0, "", "", 0, fmt.Errorf("missing or invalid tipoDte field")
		}
	}

	controlNumber, ok := identification["numeroControl"].(string)
	if !ok {
		documento := docMap["documento"].(map[string]interface{})
		controlNumber, ok = documento["numeroControl"].(string)
		if !ok {
			logs.Error("Missing or invalid numeroControl field")
			return 0, "", "", 0, fmt.Errorf("missing or invalid numeroControl field")
		}
	}

	generationCode, ok := identification["codigoGeneracion"].(string)
	if !ok {
		logs.Error("Missing or invalid codigoGeneracion field")
		return 0, "", "", 0, fmt.Errorf("missing or invalid codigoGeneracion field")
	}

	sequenceNumber, err := extractCorrelativo(controlNumber)
	if err != nil {
		logs.Error("Failed to extract correlativo", map[string]interface{}{
			"numeroControl": controlNumber,
			"error":         err.Error(),
		})
		return 0, "", "", 0, fmt.Errorf("failed to extract correlativo: %w", err)
	}

	return int(version), dteType, generationCode, sequenceNumber, nil
}

func extractCorrelativo(numeroControl string) (int, error) {
	// Separar por guiones: "DTE-03-00000000-000000000000085" -> ["DTE", "03", "00000000", "000000000000085"]
	parts := strings.Split(numeroControl, "-")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid numeroControl format")
	}
	lastPart := parts[3]

	correlativo, err := strconv.Atoi(lastPart)
	if err != nil {
		return 0, fmt.Errorf("invalid correlativo number: %w", err)
	}

	return correlativo, nil
}
