package credit_note_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

type CreditNoteSummary struct {
	*models.Summary                  // Hereda summary base
	TaxedDiscount   financial.Amount // Descuento gravado
	IVAPerception   financial.Amount // Percepción IVA 1%
	IVARetention    financial.Amount // Retención IVA 1%
	IncomeRetention financial.Amount // Retención Renta
}
