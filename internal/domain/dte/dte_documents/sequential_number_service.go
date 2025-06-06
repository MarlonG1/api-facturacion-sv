package dte_documents

import (
	"context"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"

	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type sequentialNumberService struct {
	sequentialRepo ports.SequentialNumberRepositoryPort
	authRepo       auth.AuthRepositoryPort
}

func NewSequentialNumberService(sequentialRepo ports.SequentialNumberRepositoryPort, authRepo auth.AuthRepositoryPort) SequentialNumberManager {
	return &sequentialNumberService{
		sequentialRepo: sequentialRepo,
		authRepo:       authRepo,
	}
}

func (m *sequentialNumberService) GetNextControlNumber(ctx context.Context, dteType string, branchID uint, posCode, establishmentCode *string) (string, error) {
	var controlNumber string
	defaultValue := "0000"
	if posCode == nil {
		posCode = &defaultValue
	}
	if establishmentCode == nil {
		establishmentCode = &defaultValue
	}

	// 1. Obtener el usuario por el branchID
	user, err := m.authRepo.GetByBranchID(ctx, branchID)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("SequentialNumberManager", "GetNextControlNumber", "failed to get user by branchID", err)
	}

	// 2. Obtener el siguiente número de control
	correlativeNumber, err := m.sequentialRepo.GetNext(ctx, dteType, branchID)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("SequentialNumberManager", "GetNextControlNumber", "failed to get next control number", err)
	}

	if user.YearInDTE {
		controlNumber = fmt.Sprintf("DTE-%s-%s%s-%s%011d",
			dteType,
			*establishmentCode,
			*posCode,
			utils.TimeNow().Format("2006"),
			correlativeNumber,
		)
	} else {
		controlNumber = fmt.Sprintf("DTE-%s-%s%s-%015d",
			dteType,
			*establishmentCode,
			*posCode,
			correlativeNumber,
		)
	}

	return controlNumber, nil
}
