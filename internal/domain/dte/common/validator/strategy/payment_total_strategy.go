package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/shopspring/decimal"
)

type PaymentTotalStrategy struct {
	Document interfaces.DTEDocument
}

// Validate Válida las reglas de total de pagos de un documento DTE
func (s *PaymentTotalStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || s.Document.GetSummary() == nil {
		return nil
	}

	// Validar términos de pago para crédito
	for _, payment := range s.Document.GetSummary().GetPaymentTypes() {
		if s.Document.GetSummary().GetOperationCondition() == constants.Credit {
			if payment.GetCode() == constants.BilletesMonedas {
				return dte_errors.NewDTEErrorSimple("InvalidPaymentTypeOP2")
			}

			if payment.GetTerm() == "" || payment.GetPeriod() == 0 {
				return dte_errors.NewDTEErrorSimple("InvalidPaymentTerms")
			}
		}

		if s.Document.GetSummary().GetOperationCondition() == constants.Cash {
			if payment.GetTerm() != "" || payment.GetPeriod() != 0 {
				return dte_errors.NewDTEErrorSimple("InvalidPaymentTermsOF")
			}
		}
	}

	paymentsTotal := decimal.Zero
	// Sumar todos los pagos
	for _, payment := range s.Document.GetSummary().GetPaymentTypes() {
		amount := decimal.NewFromFloat(payment.GetAmount())
		paymentsTotal = paymentsTotal.Add(amount)
	}

	operationTotal := decimal.NewFromFloat(s.Document.GetSummary().GetTotalToPay())

	// Validar que el total de pagos sea igual al total de operaciones
	if !paymentsTotal.Equal(operationTotal) {
		return dte_errors.NewDTEErrorSimple("InvalidPaymentTotal",
			paymentsTotal.InexactFloat64(), operationTotal.InexactFloat64())
	}

	return nil
}
