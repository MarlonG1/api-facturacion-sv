package constants

const (
	TaxIVA            = "20" // IVA 13%
	TaxIVAExport      = "C3" // IVA 0%
	TaxTourism        = "59" // Turismo 5%
	TaxTourismAirport = "71" // Turismo sálida del país $7.00
	TaxFOVIAL         = "D1" // FOVIAL %0.20/galón
	TaxCOTRANS        = "C8" // COTRANS $0.10/galón
	TaxSpecialOther   = "D5" // Otras tasas especiales

	TaxIvaAmount            = 0.13
	TaxIVAExportAmount      = 0.0
	TaxTourismAmount        = 0.05
	TaxTourismAirportAmount = 7.0
	TaxFOVIALAmount         = 0.20
	TaxCOTRANSAmount        = 0.10
)

var (
	// AllowedTaxTypes contiene los tipos de impuestos permitidos, usado para validaciones
	AllowedTaxTypes = []string{
		TaxIVA,
		TaxIVAExport,
		TaxTourism,
		TaxTourismAirport,
		TaxFOVIAL,
		TaxCOTRANS,
		TaxSpecialOther,
	}

	MapAllowedTaxTypes = map[string]bool{
		TaxIVA:            true,
		TaxIVAExport:      true,
		TaxTourism:        true,
		TaxTourismAirport: true,
		TaxFOVIAL:         true,
		TaxCOTRANS:        true,
		TaxSpecialOther:   true,
	}

	TaxDescriptions = map[string]string{
		TaxIVA:            "IVA 13%",
		TaxIVAExport:      "IVA de exportación 0%",
		TaxTourism:        "Turismo 5%",
		TaxTourismAirport: "Turismo sálida del país $7.00",
		TaxFOVIAL:         "FOVIAL %0.20/galón",
		TaxCOTRANS:        "COTRANS $0.10/galón",
		TaxSpecialOther:   "Otras tasas especiales",
	}
)
