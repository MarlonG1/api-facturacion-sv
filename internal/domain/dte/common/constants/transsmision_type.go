package constants

const (
	TransmisionNormal = iota + 1
	TransmisionContingencia
)

var (
	// AllowedTransmisionTypes contiene los tipos de transmisiones permitidos, usado para validaciones
	AllowedTransmisionTypes = []int{
		TransmisionNormal,
		TransmisionContingencia,
	}
)
