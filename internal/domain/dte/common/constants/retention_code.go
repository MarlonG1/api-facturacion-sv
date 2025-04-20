package constants

import "github.com/shopspring/decimal"

const (
	RetentionOnePercent      = "22" // Retención del 1% sobre el monto total
	RetentionThirteenPercent = "C4" // Retención del 13% sobre el monto total
	OtherRetentions          = "C9" // Otras retenciones (cualquier otro código de retención no especificado)
)

var (
	// AllowedRetentionCodes contiene los códigos de retención permitidos, usado para validaciones
	AllowedRetentionCodes = map[string]bool{
		RetentionOnePercent:      true,
		RetentionThirteenPercent: true,
		OtherRetentions:          true,
	}

	ListRetentionCodes = []string{
		RetentionOnePercent,
		RetentionThirteenPercent,
	}

	GetRetentionAmount = map[string]decimal.Decimal{
		RetentionOnePercent:      decimal.NewFromFloat(0.01),
		RetentionThirteenPercent: decimal.NewFromFloat(0.13),
	}
)
