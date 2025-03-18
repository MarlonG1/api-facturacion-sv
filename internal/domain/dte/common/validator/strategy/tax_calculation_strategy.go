package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type TaxCalculationStrategy struct {
	Document interfaces.DTEDocument
}

// Validate valida que los códigos de impuestos sean válidos en los items del documento
func (s *TaxCalculationStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || len(s.Document.GetItems()) == 0 {
		return nil
	}

	for _, item := range s.Document.GetItems() {
		if err := s.validateItemTaxCodes(item); err != nil {
			return err
		}
	}

	return nil
}

// validateItemTaxCodes valida que los códigos de impuestos sean válidos
func (s *TaxCalculationStrategy) validateItemTaxCodes(item interfaces.Item) *dte_errors.DTEError {
	for _, tax := range item.GetTaxes() {
		if exist := constants.MapAllowedTaxTypes[tax]; !exist {
			return dte_errors.NewDTEErrorSimple("UnsupportedTaxCode", tax)
		}
	}
	return nil
}
