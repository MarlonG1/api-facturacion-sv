package dte

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	authManager "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	dteInterfaces "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type InvalidationUseCase struct {
	dteManager          dteInterfaces.DTEManager
	authManager         authManager.AuthManager
	invalidationManager invalidation.InvalidationManager
	mapper              *request_mapper.InvalidationMapper
	transmitter         ports.BaseTransmitter
}

func NewInvalidationUseCase(dteManager dteInterfaces.DTEManager, invalidationManager invalidation.InvalidationManager, authManager authManager.AuthManager, transmitter ports.BaseTransmitter) *InvalidationUseCase {
	return &InvalidationUseCase{
		dteManager:          dteManager,
		invalidationManager: invalidationManager,
		authManager:         authManager,
		transmitter:         transmitter,
		mapper:              request_mapper.NewInvalidationMapper(),
	}
}

func (u *InvalidationUseCase) InvalidateDocument(ctx context.Context, request structs.InvalidationRequest) error {
	// 1. Sacar los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	token := ctx.Value("token").(string)

	// 2. Validar los campos del request
	if err := u.mapper.ValidateInvalidationReRequest(&request); err != nil {
		return err
	}

	// 3. Validar el estado del DTE
	if err := u.invalidationManager.ValidateStatus(ctx, claims.BranchID, request); err != nil {
		return err
	}

	// 4. Obtener el DTE Original
	originalDTE, err := u.dteManager.GetByGenerationCode(ctx, claims.BranchID, request.GenerationCode)
	if err != nil {
		return err
	}

	// 5. Obtener informacion del Issuer
	issuer, err := u.authManager.GetIssuer(ctx, claims.BranchID)
	if err != nil {
		return err
	}

	// 6. Mapear a modelo de dominio
	invalidationDocument, err := u.mapper.MapToInvalidationDocument(&request, issuer, originalDTE.Details, originalDTE.CreatedAt)
	if err != nil {
		return err
	}

	// 7. Validar el documento de invalidación
	if err = u.invalidationManager.Validate(ctx, claims.BranchID, invalidationDocument); err != nil {
		return err
	}

	// 8. Mapear a modelo de hacienda
	mhInvalidation := response_mapper.ToMHInvalidation(invalidationDocument)
	if mhInvalidation == nil {
		logs.Error("Error mapping invoice to hacienda model", map[string]interface{}{"error": "nil model"})
		return shared_error.NewGeneralServiceError("InvoiceUseCase", "InvalidateDocument", "Error mapping invoice to hacienda model", nil)
	}

	// 9. Transmitir invalidación a hacienda
	result, err := u.transmitter.RetryTransmission(ctx, mhInvalidation, token, claims.NIT)
	if err != nil {
		return err
	}
	if result.Status != ReceivedStatus {
		logs.Warn("Error transmitting invalidation", map[string]interface{}{"error": "TransmissionFailed"})
		return dte_errors.NewDTEErrorSimple("TransmissionFailed")
	}

	// 10. Invalidar documento original
	if err := u.invalidationManager.InvalidateDocument(ctx, claims.BranchID, request.GenerationCode); err != nil {
		logs.Error("Failed to update original DTE status", map[string]interface{}{
			"error": err.Error(),
			"code":  request.GenerationCode,
		})
		return err
	}

	return nil
}
