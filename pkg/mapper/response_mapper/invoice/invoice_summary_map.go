package invoice

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapInvoiceResponseSummary mapea un resumen de invoice a un modelo de resumen de invoice -> Origen: Response
func MapInvoiceResponseSummary(summary invoice_models.InvoiceSummary) *structs.InvoiceSummary {
	result := MapInvoiceSummary(summary)
	result.DescuGravada = summary.TaxedDiscount.GetValue()
	result.IvaRete1 = summary.IVARetention.GetValue()
	result.TotalIva = summary.TotalIva.GetValue()
	result.SaldoFavor = summary.BalanceInFavor.GetValue()
	result.ReteRenta = summary.IncomeRetention.GetValue()

	return result
}

func MapInvoiceSummary(summary interfaces.Summary) *structs.InvoiceSummary {
	result := &structs.InvoiceSummary{
		TotalNoSuj:          summary.GetTotalNonSubject(),
		TotalExenta:         summary.GetTotalExempt(),
		TotalGravada:        summary.GetTotalTaxed(),
		SubTotalVentas:      summary.GetSubtotalSales(),
		DescuNoSuj:          summary.GetNonSubjectDiscount(),
		DescuExenta:         summary.GetExemptDiscount(),
		DescuGravada:        summary.GetExemptDiscount(),
		PorcentajeDescuento: summary.GetDiscountPercentage(),
		TotalDescu:          summary.GetTotalDiscount(),
		SubTotal:            summary.GetSubTotal(),
		MontoTotalOperacion: summary.GetTotalOperation(),
		TotalNoGravado:      summary.GetTotalNotTaxed(),
		TotalPagar:          summary.GetTotalToPay(),
		TotalLetras:         summary.GetTotalInWords(),
		CondicionOperacion:  summary.GetOperationCondition(),
		Tributos:            common.MapTaxes(summary.GetTotalTaxes()),
	}

	// Mapear pagos
	if len(summary.GetPaymentTypes()) > 0 {
		result.Pagos = common.MapCommonResponsePayments(summary.GetPaymentTypes())
	}

	return result
}
