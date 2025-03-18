package constants

const (
	BilletesMonedas = "01" // Billetes y monedas
	TarjetaDebito   = "02" // Tarjeta de débito
	TarjetaCredito  = "03" // Tarjeta de crédito
	Cheque          = "04" // Cheque
	TransBancaria   = "05" // Transferencia bancaria
	TarjetaPrePago  = "06" // Tarjeta prepago
	Vales           = "07" // Vale
	CriptoMoneda    = "08" // Moneda virtual - Criptomoneda
	PagosElect      = "09" // Pago electrónico
	GiftCard        = "10" // Gift card
	NotaAbono       = "11" // Nota de abono
	OtraFormaPago   = "12" // Otra forma de pago
	ContoPrepago    = "13" // Cuenta prepago
	AplicaARete     = "14" // Aplicación de retención
	NoAplica        = "99" // No aplica
)

var (
	// AllowedPaymentTypes contiene los tipos de pagos permitidos, usado para validaciones
	AllowedPaymentTypes = []string{
		BilletesMonedas,
		TarjetaDebito,
		TarjetaCredito,
		Cheque,
		TransBancaria,
		TarjetaPrePago,
		Vales,
		CriptoMoneda,
		PagosElect,
		GiftCard,
		NotaAbono,
		OtraFormaPago,
		ContoPrepago,
		AplicaARete,
		NoAplica,
	}
)
