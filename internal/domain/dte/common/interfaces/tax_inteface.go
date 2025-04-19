package interfaces

// TaxGetter es una interfaz que define los métodos getter que debe implementar un impuesto
type TaxGetter interface {
	GetTotalAmount() float64 // GetTotalAmount obtiene el monto total del impuesto
	GetCode() string         // GetCode obtiene el código del impuesto
	GetDescription() string  // GetDescription obtiene la descripción del impuesto
	GetValue() float64       // GetValue obtiene el valor del impuesto
}

// TaxSetter es una interfaz que define los métodos setter que debe implementar un impuesto
type TaxSetter interface {
	SetTotalAmount(totalAmount float64) error // SetTotalAmount establece el monto total del impuesto
	SetCode(code string) error                // SetCode establece el código del impuesto
	SetDescription(description string) error  // SetDescription establece la descripción del impuesto
	SetValue(value float64) error             // SetValue establece el valor del impuesto
}

// Tax es una interfaz que combina los getters y setters de Tax
type Tax interface {
	TaxGetter
	TaxSetter
}
