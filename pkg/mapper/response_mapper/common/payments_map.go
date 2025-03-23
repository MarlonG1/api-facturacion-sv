package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// MapCommonResponsePayments mapea los pagos a un modelo de pagos -> Origen: Response
func MapCommonResponsePayments(payments []interfaces.PaymentType) []structs.DTEPayment {
	result := make([]structs.DTEPayment, len(payments))
	for i, payment := range payments {

		result[i] = structs.DTEPayment{
			Codigo:    payment.GetCode(),
			MontoPago: payment.GetAmount(),
		}

		if reference := payment.GetReference(); reference != "" {
			result[i].Referencia = utils.ToStringPointer(reference)
		}

		if term := payment.GetTerm(); term != nil {
			result[i].Plazo = utils.ToStringPointer(*term)
		}
		if period := payment.GetPeriod(); period != nil {
			result[i].Periodo = utils.ToIntPointer(*period)
		}
	}
	return result
}
