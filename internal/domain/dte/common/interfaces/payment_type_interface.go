package interfaces

// PaymentTypeGetter es una interfaz que define los métodos getter que debe implementar un tipo de pago
type PaymentTypeGetter interface {
	GetCode() string        // GetCode obtiene el código del tipo de pago
	GetAmount() float64     // GetAmount obtiene el monto del tipo de pago
	GetReference() string   // GetReference obtiene la referencia del tipo de pago
	GetTerm() *string       // GetTerm obtiene el plazo del tipo de pago
	GetPeriod() *int        // GetPeriod obtiene el periodo del tipo de pago
	GetPeriodPointer() *int // GetPeriodPointer obtiene el periodo del tipo de pago como puntero
}

// PaymentTypeSetter es una interfaz que define los métodos setter que debe implementar un tipo de pago
type PaymentTypeSetter interface {
	SetCode(code string) error           // SetCode establece el código del tipo de pago
	SetAmount(amount float64) error      // SetAmount establece el monto del tipo de pago
	SetReference(reference string) error // SetReference establece la referencia del tipo de pago
	SetTerm(term *string) error          // SetTerm establece el plazo del tipo de pago
	SetPeriod(period *int) error         // SetPeriod establece el periodo del tipo de pago
}

// PaymentType es una interfaz que combina los getters y setters de PaymentType
type PaymentType interface {
	PaymentTypeGetter
	PaymentTypeSetter
}
