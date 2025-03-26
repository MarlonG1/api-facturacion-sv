package dte

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
)

type DTEConsultUseCase struct {
	dteService interfaces.DTEManager
}

func NewDTEConsultUseCase(dteService interfaces.DTEManager) *DTEConsultUseCase {
	return &DTEConsultUseCase{
		dteService: dteService,
	}
}

func (u *DTEConsultUseCase) GetByGenerationCode(ctx context.Context, id string) (interface{}, error) {
	// 1. Obtener los claims del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)

	// 2. Consultar el documento por código de generación
	dte, err := u.dteService.GetByGenerationCode(ctx, claims.BranchID, id)
	if err != nil {
		return nil, err
	}

	return dte, nil
}
