package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/shopspring/decimal"
	"math"
)

// MapTaxes mapea los impuestos de una invoice
func MapTaxes(taxes []interfaces.Tax) []structs.DTETax {
	result := make([]structs.DTETax, 0)

	for _, tax := range taxes {
		valor := decimal.NewFromFloat(tax.GetValue())

		roundedValue := valor.Round(2)
		floatValue := roundedValue.InexactFloat64()

		if floatValue == math.Floor(floatValue) {
			floatValue = math.Round(floatValue*100) / 100
		}

		result = append(result, structs.DTETax{
			Codigo:      tax.GetCode(),
			Descripcion: tax.GetDescription(),
			Valor:       floatValue,
		})
	}

	return result
}
