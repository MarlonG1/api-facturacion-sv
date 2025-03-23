package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

// MapCommonRequestSummaryTaxes convierte un arreglo de TaxRequest a un arreglo de Tax, pero de SummaryTax, no de items
func MapCommonRequestSummaryTaxes(taxes []structs.TaxRequest) ([]interfaces.Tax, error) {
	result := make([]interfaces.Tax, len(taxes))
	for i, tax := range taxes {
		if tax.Code == "" || tax.Description == "" {
			return nil, shared_error.NewGeneralServiceError("CommonMapper", "MapCommonRequestSummaryTaxes", "Tax code and description are required, but one or both are empty", nil)
		}

		if tax.Value == 0 && tax.Code != constants.TaxIVAExport {
			return nil, dte_errors.NewValidationError("RequiredField", "Summary->Taxes->Value")
		}

		taxVO, err := financial.NewTaxType(tax.Code)
		if err != nil {
			return nil, err
		}
		taxAmount, err := financial.NewAmount(tax.Value)
		if err != nil {
			return nil, err
		}

		result[i] = &models.Tax{
			Code:        *taxVO,
			Value:       &models.TaxAmount{TotalAmount: *taxAmount},
			Description: tax.Description,
		}
	}

	return result, nil
}
