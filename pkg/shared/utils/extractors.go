package utils

import (
	"encoding/json"
)

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

type AuxiliarSummaryExtractor struct {
	Summary struct {
		SubTotal     float64 `json:"subTotal"`
		IvaRetention float64 `json:"ivaRete1"`
	} `json:"resumen"`
}

type AuxiliarTotalAmountsExtractor struct {
	Summary struct {
		TotalTaxed      float64 `json:"totalGravada"`
		TotalExempt     float64 `json:"totalExenta"`
		TotalNotSubject float64 `json:"totalNoSuj"`
	} `json:"resumen"`
}

type AuxiliarDTETotalAmountsExtractor struct {
	Summary struct {
		TotalTaxed      float64 `json:"total_taxed"`
		TotalExempt     float64 `json:"total_exempt"`
		TotalNotSubject float64 `json:"total_non_subject"`
	} `json:"summary"`
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

func ExtractAuxiliarSummaryFromStringJSON(document interface{}) (AuxiliarSummaryExtractor, error) {
	var summary AuxiliarSummaryExtractor

	// 1. Convertir a formato JSON el documento
	jsonData := []byte(document.(string))

	// 2. Extraer parte de la información del Resumen
	if err := json.Unmarshal(jsonData, &summary); err != nil {
		return summary, err
	}

	return summary, nil
}

func ExtractSummaryTotalAmounts(document interface{}) (AuxiliarTotalAmountsExtractor, error) {
	var summary AuxiliarTotalAmountsExtractor

	// 1. Convertir a formato JSON el documento
	jsonData, err := json.Marshal(document)
	if err != nil {
		return summary, err
	}

	// 2. Extraer parte de la información del Resumen
	if err = json.Unmarshal(jsonData, &summary); err != nil {
		return summary, err
	}

	return summary, nil
}

func ExtractSummaryTotalAmountsFromStringJSON(document interface{}) (AuxiliarTotalAmountsExtractor, error) {
	var summary AuxiliarTotalAmountsExtractor

	// 1. Convertir a formato JSON el documento
	jsonData := []byte(document.(string))

	// 2. Extraer parte de la información del Resumen
	if err := json.Unmarshal(jsonData, &summary); err != nil {
		return summary, err
	}

	return summary, nil
}
