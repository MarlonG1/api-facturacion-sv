package invoice_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

type InvoiceSummary struct {
	*models.Summary
	TaxedDiscount   financial.Amount `json:"taxedDiscount"`   // descuGravada
	IVAPerception   financial.Amount `json:"IVAPerception"`   // ivaPerci1
	IVARetention    financial.Amount `json:"IVARetention"`    // ivaRete1
	IncomeRetention financial.Amount `json:"incomeRetention"` // reteRenta
	TotalIva        financial.Amount `json:"totalIva"`        // totalIva
	BalanceInFavor  financial.Amount `json:"balanceInFavor"`  // saldoFavor
}
