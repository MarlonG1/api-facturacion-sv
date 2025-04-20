package interfaces

// Summary es una interfaz que define los métodos que debe implementar un resumen
type Summary interface {
	SummaryGetters
	SummarySetters
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

// SummarySetters es una interfaz que define los métodos setter que debe implementar un resumen
type SummarySetters interface {
	SetTotalNonSubject(totalNonSubject float64) error       // SetTotalNonSubject establece el total de los items no sujetos
	SetTotalExempt(totalExempt float64) error               // SetTotalExempt establece el total de los items exentos
	SetTotalTaxed(totalTaxed float64) error                 // SetTotalTaxed establece el total de los items gravados
	SetSubTotal(subTotal float64) error                     // SetSubTotal establece el subtotal
	SetSubtotalSales(subtotalSales float64) error           // SetSubtotalSales establece el subtotal de ventas
	SetNonSubjectDiscount(nonSubjectDiscount float64) error // SetNonSubjectDiscount establece el descuento de los items no sujetos
	SetExemptDiscount(exemptDiscount float64) error         // SetExemptDiscount establece el descuento de los items exentos
	SetDiscountPercentage(discountPercentage float64) error // SetDiscountPercentage establece el porcentaje de descuento
	SetTotalDiscount(totalDiscount float64) error           // SetTotalDiscount establece el total de descuentos
	SetTotalTaxes(totalTaxes []Tax) error                   // SetTotalTaxes establece los impuestos totales
	SetTotalOperation(totalOperation float64) error         // SetTotalOperation establece el total de operación
	SetTotalNotTaxed(totalNotTaxed float64) error           // SetTotalNotTaxed establece el total no gravado
	SetPaymentTypes(paymentTypes []PaymentType) error       // SetPaymentTypes establece los tipos de pago
	SetOperationCondition(operationCondition int) error     // SetOperationCondition establece la condición de pago
	SetElectronicPayment(electronicPayment *string) error   // SetElectronicPayment establece el medio de pago electrónico
	SetTotalInWords(totalInWords string) error              // SetTotalInWords establece el total en palabras
	SetTotalToPay(totalToPay float64) error                 // SetTotalToPay establece el total a pagar
	SetForceTotalToPay(totalToPay float64)                  // SetForceTotalToPay establece el total a pagar sin validación
}

// SummaryManager es una interfaz que combina los getters y setters de Summary
type SummaryManager interface {
	SummaryGetters
	SummarySetters
}
