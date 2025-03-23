package common

// MapTaxCodes mapea los códigos de impuestos
func MapTaxCodes(taxes []string) {
	codes := make([]string, len(taxes))
	for i, tax := range taxes {
		codes[i] = tax
	}
}
