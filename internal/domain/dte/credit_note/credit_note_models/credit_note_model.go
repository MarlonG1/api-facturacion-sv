package credit_note_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type CreditNoteModel struct {
	*models.DTEDocument
	CreditItems   []CreditNoteItem
	CreditSummary CreditNoteSummary
}
