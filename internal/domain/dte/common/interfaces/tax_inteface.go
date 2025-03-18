package interfaces

// Tax es una interfaz que define los métodos que debe implementar un impuesto
type Tax interface {
	GetTotalAmount() float64 // GetTotalAmount obtiene el monto total del impuesto
	GetCode() string         // GetCode obtiene el código del impuesto
	GetDescription() string  // GetDescription obtiene la descripción del impuesto
	GetValue() float64       // GetValue obtiene el valor del impuesto
}
