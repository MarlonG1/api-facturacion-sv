package interfaces

// Item es una interfaz que define los métodos que debe implementar un item de un documento tributario electrónico
type Item interface {
	GetQuantity() float64   // GetQuantity retorna la cantidad del item
	GetItemCode() string    // GetItemCode retorna el número del item
	GetDescription() string // GetDescription retorna la descripción del item
	GetType() int           /// GetType retorna el tipo de item
	GetUnitPrice() float64  // GetUnitPrice retorna el precio unitario del item
	GetDiscount() float64   // GetDiscount retorna el descuento aplicado al item
	GetTaxes() []string     // GetTaxes retorna los impuestos aplicados al item
	GetRelatedDoc() *string // GetRelatedDoc retorna el documento relacionado al item
	GetNumber() int         // GetNumber retorna el número del item
	GetUnitMeasure() int    // GetUnitMeasure retorna la unidad de medida del item
}
