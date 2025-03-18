package constants

const (
	Producto          = iota + 1 //Tipo de item Producto
	Servicio                     //Tipo de item Servicio
	ProductoYServicio            //Tipo de item Producto y Servicio
	Impuesto                     //Tipo de item Impuesto
)

var (
	// AllowedItemTypes contiene los tipos de items permitidos para validaciones de tipo.
	// Estos valores representan los distintos tipos de items en el sistema: Producto, Servicio, Ambos e Impuesto.
	AllowedItemTypes = []int{
		Producto,
		Servicio,
		ProductoYServicio,
		Impuesto,
	}
)
