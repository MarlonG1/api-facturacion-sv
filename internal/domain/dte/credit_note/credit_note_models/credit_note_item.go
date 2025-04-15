package credit_note_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
)

type CreditNoteItem struct {
	*models.Item
	NonSubjectSale financial.Amount
	ExemptSale     financial.Amount
	TaxedSale      financial.Amount
	SuggestedPrice financial.Amount
	NonTaxed       financial.Amount
}
