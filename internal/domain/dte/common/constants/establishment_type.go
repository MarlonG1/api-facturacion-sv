package constants

const (
	Sucursal       = "01" // Sucursal o Agencia
	CasaMatriz     = "02" // Casa Matriz
	DepositoBodega = "04" // Depósito o Bodega para distribución
	PredioOPatio   = "07" // Predio o Patio para distribución
	Otro           = "20" // Otro tipo no especificado
)

var (
	// AllowedEstablishmentTypes contiene los tipos de establecimientos permitidos
	AllowedEstablishmentTypes = []string{
		CasaMatriz,
		Sucursal,
		DepositoBodega,
		PredioOPatio,
		Otro,
	}
)
