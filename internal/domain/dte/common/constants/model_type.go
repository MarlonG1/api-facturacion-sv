package constants

const (
	ModeloFacturacionPrevio = iota + 1
	ModeloFacturacionDiferido
)

var (
	// AllowedModeloFacturacion contiene los tipos de modelos de facturación permitidos, usado para validaciones
	AllowedModeloFacturacion = []int{
		ModeloFacturacionPrevio,
		ModeloFacturacionDiferido,
	}
)
