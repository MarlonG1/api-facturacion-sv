package constants

const (
	Cash   = iota + 1 // Representa una condición de pago en efectivo
	Credit            // Representa una condición de pago a crédito
	Other             // Representa una condición de pago diferente a efectivo o crédito
)

var (
	// ValidPaymentConditions Es una lista de valores permitidos para el campo OperationCondition
	ValidPaymentConditions = []int{
		Cash,
		Credit,
		Other,
	}
)
