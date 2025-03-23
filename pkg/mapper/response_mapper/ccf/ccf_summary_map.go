package ccf

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func MapCCFResponseSummary(summary ccf_models.CreditSummary) *structs.DTESummary {
	ivaPerci1 := summary.IVAPerception.GetValue()
	result := common.MapCommonResponseSummary(summary)
	result.DescuGravada = summary.TaxedDiscount.GetValue()
	result.IvaRete1 = summary.IVARetention.GetValue()
	result.IvaPerci1 = &ivaPerci1
	result.ReteRenta = summary.IncomeRetention.GetValue()
	result.SaldoFavor = summary.BalanceInFavor.GetValue()
	return result
}
