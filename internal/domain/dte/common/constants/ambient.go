package constants

/*
AllowedAmbientValues Es una lista de valores permitidos para el campo Ambient
"00" - Testing
"01" - Production
*/
const (
	Testing    = "00"
	Production = "01"
)

var (
	// AllowedAmbientValues Es una lista de valores permitidos para el campo Ambient, usado para validaciones
	AllowedAmbientValues = []string{
		Testing,
		Production,
	}
)
