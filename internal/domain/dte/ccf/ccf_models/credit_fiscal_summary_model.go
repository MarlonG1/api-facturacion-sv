package ccf_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

type CreditSummary struct {
	*models.Summary                          // Hereda summary base
	TaxedDiscount           financial.Amount // Descuento gravado
	IVAPerception           financial.Amount // Percepción IVA 1%
	IVARetention            financial.Amount // Retención IVA 1%
	BalanceInFavor          financial.Amount // Saldo a Favor
	IncomeRetention         financial.Amount // Retención Renta
	ElectronicPaymentNumber *string          // Número de pago electrónico
}
