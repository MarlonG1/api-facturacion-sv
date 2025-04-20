package interfaces

// ItemGetter es una interfaz que define los métodos getter que debe implementar un item
type ItemGetter interface {
	GetQuantity() float64   // GetQuantity retorna la cantidad del item
	GetItemCode() string    // GetItemCode retorna el número del item
	GetDescription() string // GetDescription retorna la descripción del item
	GetType() int           // GetType retorna el tipo de item
	GetUnitPrice() float64  // GetUnitPrice retorna el precio unitario del item
	GetDiscount() float64   // GetDiscount retorna el descuento aplicado al item
	GetTaxes() []string     // GetTaxes retorna los impuestos aplicados al item
	GetRelatedDoc() *string // GetRelatedDoc retorna el documento relacionado al item
	GetNumber() int         // GetNumber retorna el número del item
	GetUnitMeasure() int    // GetUnitMeasure retorna la unidad de medida del item
}

// ItemSetter es una interfaz que define los métodos setter que debe implementar un item
type ItemSetter interface {
	SetQuantity(quantity float64) error      // SetQuantity establece la cantidad del item
	SetItemCode(itemCode string) error       // SetItemCode establece el número del item
	SetDescription(description string) error // SetDescription establece la descripción del item
	SetType(itemType int) error              // SetType establece el tipo de item
	SetUnitPrice(unitPrice float64) error    // SetUnitPrice establece el precio unitario del item
	SetForceUnitPrice(unitPrice float64)     // SetUnitPrice establece el precio unitario del item
	SetDiscount(discount float64) error      // SetDiscount establece el descuento aplicado al item
	SetTaxes(taxes []string) error           // SetTaxes establece los impuestos aplicados al item
	SetRelatedDoc(relatedDoc *string) error  // SetRelatedDoc establece el documento relacionado al item
	SetForceRelatedDoc(relatedDoc *string)   // SetRelatedDoc establece el documento relacionado al item
	SetNumber(number int) error              // SetNumber establece el número del item
	SetUnitMeasure(unitMeasure int) error    // SetUnitMeasure establece la unidad de medida del item
}

// Item es una interfaz que combina los getters y setters de Item
type Item interface {
	ItemGetter
	ItemSetter
}
