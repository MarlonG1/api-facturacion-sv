package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseSummary mapea un resumen de invoice a un modelo de resumen de invoice -> Origen: Response
func MapCommonResponseSummary(summary interfaces.Summary) *structs.DTESummary {
	result := &structs.DTESummary{
		TotalNoSuj:          summary.GetTotalNonSubject(),
		TotalExenta:         summary.GetTotalExempt(),
		TotalGravada:        summary.GetTotalTaxed(),
		SubTotalVentas:      summary.GetSubtotalSales(),
		DescuNoSuj:          summary.GetNonSubjectDiscount(),
		DescuExenta:         summary.GetExemptDiscount(),
		PorcentajeDescuento: summary.GetDiscountPercentage(),
		TotalDescu:          summary.GetTotalDiscount(),
		SubTotal:            summary.GetSubTotal(),
		MontoTotalOperacion: summary.GetTotalOperation(),
		TotalNoGravado:      summary.GetTotalNotTaxed(),
		TotalPagar:          summary.GetTotalToPay(),
		TotalLetras:         summary.GetTotalInWords(),
		CondicionOperacion:  summary.GetOperationCondition(),
		Tributos:            MapTaxes(summary.GetTotalTaxes()),
	}

	// Mapear pagos
	if len(summary.GetPaymentTypes()) > 0 {
		result.Pagos = MapCommonResponsePayments(summary.GetPaymentTypes())
	}

	return result
}
