package utils

import "encoding/json"

type AuxiliarIdentificationExtractor struct {
	Identification struct {
		DTEType        string `json:"tipoDte"`
		ControlNumber  string `json:"numeroControl"`
		GenerationCode string `json:"codigoGeneracion"`
	} `json:"identificacion"`
	Issuer struct {
		NIT string `json:"nit"`
	} `json:"emisor"`
}

func ExtractAuxiliarIdentification(document interface{}) (AuxiliarIdentificationExtractor, error) {
	var identification AuxiliarIdentificationExtractor

	// 1. Convertir a formato JSON el documento
	jsonData, err := json.Marshal(document)
	if err != nil {
		return identification, err
	}

	// 2. Extraer parte de la información de Identificación
	if err := json.Unmarshal(jsonData, &identification); err != nil {
		return identification, err
	}

	return identification, nil
}

func ExtractAuxiliarIdentificationFromStringJSON(document interface{}) (AuxiliarIdentificationExtractor, error) {
	var identification AuxiliarIdentificationExtractor

	// 1. Convertir a formato JSON el documento
	jsonData := []byte(document.(string))

	// 2. Extraer parte de la información de Identificación
	if err := json.Unmarshal(jsonData, &identification); err != nil {
		return identification, err
	}

	return identification, nil
}

func ExtractAuxiliarDTEInfo(document interface{}) (AuxiliarIdentificationExtractor, error) {
	var dteInfo AuxiliarIdentificationExtractor

	// 1. Convertir a formato JSON el documento
	jsonData, err := json.Marshal(document)
	if err != nil {
		return dteInfo, err
	}

	// 2. Extraer parte de la información del DTE
	if err := json.Unmarshal(jsonData, &dteInfo); err != nil {
		return dteInfo, err
	}

	return dteInfo, nil
}
