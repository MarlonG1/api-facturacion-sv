package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// MapCommonRequestSummary mapea un resumen común a un modelo de resumen -> Origen: Request
func MapCommonRequestSummary(summary structs.SummaryRequest) (*models.Summary, error) {
	if err := validateSummaryFields(summary); err != nil {
		return nil, err
	}

	totalNonSubject, err := financial.NewAmountForTotal(summary.TotalNonSubject)
	if err != nil {
		return nil, err
	}

	totalExempt, err := financial.NewAmountForTotal(summary.TotalExempt)
	if err != nil {
		return nil, err
	}

	totalTaxed, err := financial.NewAmountForTotal(summary.TotalTaxed)
	if err != nil {
		return nil, err
	}

	subTotal, err := financial.NewAmountForTotal(summary.SubTotal)
	if err != nil {
		return nil, err
	}

	nonSubjectDiscount, err := financial.NewAmountForTotal(summary.NonSubjectDiscount)
	if err != nil {
		return nil, err
	}

	exemptDiscount, err := financial.NewAmountForTotal(summary.ExemptDiscount)
	if err != nil {
		return nil, err
	}

	discountPercentage, err := financial.NewDiscount(summary.DiscountPercentage)
	if err != nil {
		return nil, err
	}

	totalDiscount, err := financial.NewAmountForTotal(summary.TotalDiscount)
	if err != nil {
		return nil, err
	}

	subTotalSales, err := financial.NewAmountForTotal(summary.SubTotalSales)
	if err != nil {
		return nil, err
	}

	totalOperation, err := financial.NewAmountForTotal(summary.TotalOperation)
	if err != nil {
		return nil, err
	}

	totalNonTaxed, err := financial.NewAmountForTotal(summary.TotalNonTaxed)
	if err != nil {
		return nil, err
	}

	paymentCondition, err := financial.NewPaymentCondition(summary.OperationCondition)
	if err != nil {
		return nil, err
	}

	totalToPay, err := financial.NewAmountForTotal(summary.TotalToPay)
	if err != nil {
		return nil, err
	}

	taxes, err := MapCommonRequestSummaryTaxes(summary.Taxes)
	if err != nil {
		return nil, err
	}

	paymentTypes, err := MapCommonRequestPaymentsType(summary.PaymentTypes)
	if err != nil {
		return nil, err
	}

	return &models.Summary{
		TotalNonSubject:    *totalNonSubject,
		TotalExempt:        *totalExempt,
		TotalTaxed:         *totalTaxed,
		SubTotal:           *subTotal,
		NonSubjectDiscount: *nonSubjectDiscount,
		ExemptDiscount:     *exemptDiscount,
		DiscountPercentage: *discountPercentage,
		TotalDiscount:      *totalDiscount,
		SubTotalSales:      *subTotalSales,
		TotalOperation:     *totalOperation,
		TotalNonTaxed:      *totalNonTaxed,
		OperationCondition: *paymentCondition,
		TotalToPay:         *totalToPay,
		TotalTaxes:         taxes,
		PaymentTypes:       paymentTypes,
		TotalInWords:       *summary.TotalInWords,
	}, nil
}

func validateSummaryFields(summary structs.SummaryRequest) error {
	//if summary.TotalTaxed == 0 {
	//	return dte_errors.NewValidationError("RequiredField", "TotalTaxed")
	//}
	if summary.SubTotal == 0 {
		return dte_errors.NewValidationError("RequiredField", "SubTotal")
	}
	if summary.SubTotalSales == 0 {
		return dte_errors.NewValidationError("RequiredField", "SubTotalSales")
	}
	if summary.TotalOperation == 0 {
		return dte_errors.NewValidationError("RequiredField", "TotalOperation")
	}
	if summary.TotalToPay == 0 {
		return dte_errors.NewValidationError("RequiredField", "TotalToPay")
	}

	return nil
}
