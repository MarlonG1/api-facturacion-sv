package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// MapCommonRequestPaymentsType mapea un arreglo de pagos a un modelo de pagos -> Origen: Request
func MapCommonRequestPaymentsType(payments []structs.PaymentRequest) ([]interfaces.PaymentType, error) {
	result := make([]interfaces.PaymentType, len(payments))

	for i, payment := range payments {
		var term *financial.PaymentTerm

		if payment.Code == "" || payment.Amount == 0 {
			return nil, shared_error.NewGeneralServiceError(
				"CommonMapper",
				"MapCommonRequestPaymentsType",
				"Code and Amount are required and cannot be empty in payment_types",
				nil,
			)
		}

		paymentCode, err := financial.NewPaymentType(payment.Code)
		if err != nil {
			return nil, err
		}

		amount, err := financial.NewAmount(payment.Amount)
		if err != nil {
			return nil, err
		}

		if payment.Term != nil {
			term, err = financial.NewPaymentTerm(*payment.Term)
			if err != nil {
				return nil, err
			}
		}

		if payment.Reference == nil {
			payment.Reference = new(string)
		}

		result[i] = &models.PaymentType{
			Code:      *paymentCode,
			Amount:    *amount,
			Reference: *payment.Reference,
			Period:    payment.Period,
			Term:      term,
		}
	}

	return result, nil
}
