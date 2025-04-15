package invalidation

import (
	"context"
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/validator"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type invalidationService struct {
	validator  *validator.InvalidationRulesValidator
	dteManager dte_documents.DTEManager
}

// NewInvalidationService crea una nueva instancia de InvalidationManager
func NewInvalidationService(dteManager dte_documents.DTEManager) InvalidationManager {
	return &invalidationService{
		validator:  validator.NewInvalidationRulesValidator(nil),
		dteManager: dteManager,
	}
}

func (s *invalidationService) InvalidateDocument(ctx context.Context, branchID uint, originalCode string) error {
	return s.dteManager.UpdateDTE(ctx, branchID, dte.DTEDetails{
		ID:     originalCode,
		Status: constants.DocumentInvalid,
	})
}

func (s *invalidationService) Validate(ctx context.Context, branchID uint, document *models.InvalidationDocument) error {
	// 1. Validar el documento de invalidación
	s.validator = validator.NewInvalidationRulesValidator(document)
	if err := s.validator.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *invalidationService) ValidateStatus(ctx context.Context, branchID uint, req structs.CreateInvalidationRequest) error {
	// 1. Validar la existencia del DTE y que no esté invalidado o rechazado
	if err := s.validateDTEStatus(ctx,
		branchID,
		req.GenerationCode,
		"document to invalidate",
	); err != nil {
		return err
	}

	// 2. Si la invalidacion es tipo 2, verificar que el DTE de reemplazo esté en estado correcto
	if req.Reason.Type != 2 && req.ReplacementGenerationCode != nil {
		if err := s.validateDTEStatus(ctx,
			branchID,
			*req.ReplacementGenerationCode,
			"replacement document",
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *invalidationService) validateDTEStatus(ctx context.Context, branchID uint, originalCode, message string) error {
	status, err := s.dteManager.VerifyStatus(ctx, branchID, originalCode)
	if err != nil {
		return err
	}
	return s.handleError(status, message)
}

func (s *invalidationService) handleError(status, message string) error {
	switch status {
	case constants.DocumentInvalid:
		return shared_error.NewGeneralServiceError("InvalidationService", "InvalidateDocument", fmt.Sprintf("the %s is already invalid", message), nil)
	case constants.DocumentRejected:
		return shared_error.NewGeneralServiceError("InvalidationService", "InvalidateDocument", fmt.Sprintf("the %s has been rejected by Hacienda, therefore it cannot be invalidated", message), nil)
	case constants.DocumentPending:
		return shared_error.NewGeneralServiceError("InvalidationService", "InvalidateDocument", fmt.Sprintf("the %s is still pending, therefore it cannot be invalidated", message), nil)
	default:
		return nil
	}
}
