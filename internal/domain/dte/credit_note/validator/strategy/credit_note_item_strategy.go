package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type CreditNoteItemStrategy struct {
	Document *credit_note_models.CreditNoteModel
}

func (s *CreditNoteItemStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || len(s.Document.CreditItems) == 0 {
		return dte_errors.NewDTEErrorSimple("RequiredField", "CreditItems")
	}

	// Validar número máximo de ítems
	if len(s.Document.CreditItems) > 2000 {
		return dte_errors.NewDTEErrorSimple("ExceededItemsLimit", len(s.Document.CreditItems))
	}

	for _, item := range s.Document.CreditItems {
		// Validar tipos de venta y sus restricciones
		if err := s.validateItemSaleTypes(&item); err != nil {
			return err
		}

		// Validar reglas específicas de cada ítem
		if err := s.validateItem(&item); err != nil {
			return err
		}

		// Validar reglas específicas de tipo 4 para Nota de Crédito
		if err := s.validateCreditNoteType4Rules(&item); err != nil {
			return err
		}
	}

	return nil
}

func (s *CreditNoteItemStrategy) validateItem(item *credit_note_models.CreditNoteItem) *dte_errors.DTEError {
	if item.TaxedSale.GetValue() > 0 && item.GetUnitPrice() == 0 {
		logs.Error("Unit price cannot be zero when taxed sale is present", map[string]interface{}{
			"itemNumber": item.GetNumber(),
			"taxedSale":  item.TaxedSale.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("InvalidUnitPriceZero",
			item.GetNumber(), item.GetUnitPrice(), item.TaxedSale.GetValue())
	}

	// Validación específica para Nota de Crédito: los ítems deben tener documentos relacionados
	if item.GetRelatedDoc() == nil {
		logs.Error("Related document is required for credit note items", map[string]interface{}{
			"itemNumber": item.GetNumber(),
		})
		return dte_errors.NewDTEErrorSimple("MissingItemRelatedDoc", item.GetNumber())
	}

	// Validación de impuestos: al menos uno debe estar presente si hay venta gravada
	if item.TaxedSale.GetValue() > 0 {
		if item.GetTaxes() == nil || len(item.GetTaxes()) == 0 {
			logs.Error("At least one tax is required for credit note items", map[string]interface{}{
				"itemNumber": item.GetNumber(),
			})
			return dte_errors.NewDTEErrorSimple("MissingItemTaxes", item.GetNumber())
		}
	}

	if item.GetType() != constants.Impuesto {
		for _, tax := range item.GetTaxes() {
			if !constants.MapAllowedTaxTypes[tax] {
				logs.Error("Invalid tax type", map[string]interface{}{
					"itemNumber": item.GetNumber(),
					"tax":        tax,
				})
				return dte_errors.NewDTEErrorSimple("InvalidTaxType", item.GetNumber(), tax)
			}
		}
	}

	return nil
}

func (s *CreditNoteItemStrategy) validateItemSaleTypes(item *credit_note_models.CreditNoteItem) *dte_errors.DTEError {
	// Validar que no haya ventas mixtas
	salesTypes := 0
	if item.TaxedSale.GetValue() > 0 {
		salesTypes++
	}
	if item.ExemptSale.GetValue() > 0 {
		salesTypes++
	}
	if item.NonSubjectSale.GetValue() > 0 {
		salesTypes++
	}

	if salesTypes > 1 {
		logs.Error("Mixed sales types in single item", map[string]interface{}{
			"itemNumber":     item.GetNumber(),
			"taxedSale":      item.TaxedSale.GetValue(),
			"exemptSale":     item.ExemptSale.GetValue(),
			"nonSubjectSale": item.NonSubjectSale.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("MixedSalesTypesNotAllowed", item.GetNumber())
	}

	return nil
}

func (s *CreditNoteItemStrategy) validateCreditNoteType4Rules(item interfaces.Item) *dte_errors.DTEError {
	if item.GetType() == constants.Impuesto {
		// Para tipo 4 en Nota de Crédito, validar que:
		// 1. Unidad de medida sea 99
		if item.GetUnitMeasure() != 99 {
			return dte_errors.NewDTEErrorSimple("InvalidUnitMeasure", item.GetUnitMeasure())
		}

		// 2. Solo tenga el impuesto IVA (20)
		if len(item.GetTaxes()) != 1 || item.GetTaxes()[0] != constants.TaxIVA {
			return dte_errors.NewDTEErrorSimple("InvalidTaxRulesCCF")
		}
	}

	return nil
}
