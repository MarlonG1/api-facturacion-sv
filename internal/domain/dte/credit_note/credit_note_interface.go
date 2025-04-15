package credit_note

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"golang.org/x/net/context"
)

// CreditNoteManager es la interfaz que define las operaciones de un manager de Nota de Crédito
type CreditNoteManager interface {
	// Create crea una invoice electrónica a partir de los datos de una Nota de Crédito
	Create(ctx context.Context, data *credit_note_models.CreditNoteInput, branchID uint) (*credit_note_models.CreditNoteModel, error)
	// Validate valida una nota de crédito electrónica
	Validate(creditNote *credit_note_models.CreditNoteModel) error
	// IsValid verifica si una nota de crédito electrónica es válida
	IsValid(creditNote *credit_note_models.CreditNoteModel) bool
}
