package interfaces

// Summary es una interfaz que define los métodos que debe implementar un resumen
type Summary interface {
	SummaryGetters
}

type SummaryGetters interface {
	GetTotalNonSubject() float64    // GetTotalNonSubject retorna el total de los items no sujetos
	GetTotalExempt() float64        // GetTotalExempt retorna el total de los items exentos
	GetTotalTaxed() float64         // GetTotalTaxed retorna el total de los items gravados
	GetSubTotal() float64           // GetSubTotal retorna el subtotal
	GetSubtotalSales() float64      // GetSubtotalSales retorna el subtotal de ventas
	GetNonSubjectDiscount() float64 // GetNonSubjectDiscount retorna el descuento de los items no sujetos
	GetExemptDiscount() float64     // GetExemptDiscount retorna el descuento de los items exentos
	GetDiscountPercentage() float64 // GetDiscountPercentage retorna el porcentaje de descuento
	GetTotalDiscount() float64      // GetTotalDiscount retorna el total de descuentos
	GetTotalTaxes() []Tax           // GetTotalTaxes retorna los impuestos totales
	GetTotalOperation() float64     // GetTotalOperation retorna el total de operación
	GetTotalNotTaxed() float64      // GetTotalNotTaxed retorna el total no gravado
	GetPaymentTypes() []PaymentType // GetPaymentTypes retorna los tipos de pago
	GetOperationCondition() int     // GetPaymentCondition retorna la condición de pago
	GetElectronicPayment() *string  // GetElectronicPayment retorna el medio de pago electrónico
	GetTotalInWords() string        // GetTotalInWords retorna el total en palabras
	GetTotalToPay() float64         // GetTotalToPay retorna el total a pagar
}
