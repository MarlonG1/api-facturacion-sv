package dte

import (
	"context"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

// DTEOperations define operaciones específicas para cada tipo de DTE
type DTEOperations struct{}

type AdditionalOperationsFunc func(ctx context.Context, result interface{}, branchID uint, mhModel interface{}) error

// NewDTEOperations crea una nueva instancia de DTEOperations
func NewDTEOperations() *DTEOperations {
	return &DTEOperations{}
}

// GetCreditNoteOperations devuelve las operaciones adicionales para notas de crédito
func (o *DTEOperations) GetCreditNoteOperations(dteService dte_documents.DTEManager) AdditionalOperationsFunc {
	return func(ctx context.Context, result interface{}, branchID uint, mhModel interface{}) error {
		creditNote, ok := result.(*credit_note_models.CreditNoteModel)
		if !ok {
			return fmt.Errorf("invalid result type")
		}

		for _, relatedDoc := range creditNote.GetRelatedDocuments() {
			if relatedDoc.GetGenerationType() == constants.ElectronicDocument {
				err := dteService.GenerateBalanceTransaction(
					ctx,
					branchID,
					constants.NotaCreditoElectronica,
					relatedDoc.GetDocumentNumber(),
					creditNote.GetIdentification().GetGenerationCode(),
					mhModel,
				)
				if err != nil {
					logs.Warn("Failed to generate balance transaction", map[string]interface{}{"error": err.Error()})
					return err
				}
			}
		}
		return nil
	}
}

// GetNoOperation devuelve una función vacía para DTEs sin operaciones adicionales
func (o *DTEOperations) GetNoOperation() AdditionalOperationsFunc {
	return func(ctx context.Context, result interface{}, branchID uint, mhModel interface{}) error {
		return nil
	}
}
