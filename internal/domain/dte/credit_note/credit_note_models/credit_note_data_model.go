package credit_note_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type CreditNoteInput struct {
	*models.InputDataCommon
	Items         []CreditNoteItem
	CreditSummary *CreditNoteSummary
}
