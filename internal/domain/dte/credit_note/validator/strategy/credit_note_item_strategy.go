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

		// Validar reglas específicas de no gravados
		if err := s.validateItemNonTaxedRules(&item); err != nil {
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

	return nil
}

func (s *CreditNoteItemStrategy) validateItemNonTaxedRules(item *credit_note_models.CreditNoteItem) *dte_errors.DTEError {
	nonTaxed := item.NonTaxed.GetValue()
	if nonTaxed > 0 {
		if item.GetUnitPrice() != 0 {
			logs.Error("Unit price must be zero when non_taxed > 0", map[string]interface{}{
				"itemNumber": item.GetNumber(),
				"unitPrice":  item.GetUnitPrice(),
				"nonTaxed":   nonTaxed,
			})
			return dte_errors.NewDTEErrorSimple("InvalidUnitPriceForNonTaxed",
				item.GetNumber())
		}
	}
	return nil
}

func (s *CreditNoteItemStrategy) validateItemSaleTypes(item *credit_note_models.CreditNoteItem) *dte_errors.DTEError {
	// Validar venta no gravada
	if item.NonTaxed.GetValue() > 0 {
		// No debe tener otros tipos de venta
		if item.TaxedSale.GetValue() > 0 || item.ExemptSale.GetValue() > 0 || item.NonSubjectSale.GetValue() > 0 {
			logs.Error("Items with non-taxed amount cannot have other sale types", map[string]interface{}{
				"itemNumber":     item.GetNumber(),
				"nonTaxed":       item.NonTaxed.GetValue(),
				"taxedSale":      item.TaxedSale.GetValue(),
				"exemptSale":     item.ExemptSale.GetValue(),
				"nonSubjectSale": item.NonSubjectSale.GetValue(),
			})
			return dte_errors.NewDTEErrorSimple("InvalidMixedSalesWithNonTaxed", item.GetNumber())
		}

		// No debe tener impuestos
		if len(item.GetTaxes()) > 0 {
			logs.Error("Items with non-taxed amount cannot have taxes", map[string]interface{}{
				"itemNumber": item.GetNumber(),
				"taxes":      item.GetTaxes(),
			})
			return dte_errors.NewDTEErrorSimple("InvalidTaxesWithNonTaxed", item.GetNumber())
		}
	}

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
	if item.NonTaxed.GetValue() > 0 {
		salesTypes++
	}

	if salesTypes > 1 {
		logs.Error("Mixed sales types in single item", map[string]interface{}{
			"itemNumber":     item.GetNumber(),
			"taxedSale":      item.TaxedSale.GetValue(),
			"exemptSale":     item.ExemptSale.GetValue(),
			"nonSubjectSale": item.NonSubjectSale.GetValue(),
			"nonTaxed":       item.NonTaxed.GetValue(),
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
