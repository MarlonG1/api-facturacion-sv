package interfaces

// PaymentType es una interfaz que define los métodos que debe implementar un tipo de pago
type PaymentType interface {
	GetCode() string        // GetCode obtiene el código del tipo de pago
	GetAmount() float64     // GetAmount obtiene el monto del tipo de pago
	GetReference() string   // GetReference obtiene la referencia del tipo de pago
	GetTerm() string        // GetTerm obtiene el plazo del tipo de pago
	GetPeriod() int         // GetPeriod obtiene el periodo del tipo de pago
	GetPeriodPointer() *int // GetPeriodPointer obtiene el periodo del tipo de pago como puntero
}
