package dte

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	dte, err := u.dteService.GetByGenerationCodeConsult(ctx, claims.BranchID, id)
	if err != nil {
		return nil, err
	}

	return dte, nil
}

func (u *DTEConsultUseCase) GetAllDTEs(ctx context.Context, r *http.Request) (*dte.DTEListResponse, error) {
	// 1. Parsear los parámetros de consulta
	filters, err := parseDTEFilters(r)
	if err != nil {
		return nil, err
	}

	// 2. Obtener todos los DTEs
	response, err := u.dteService.GetAllDTEs(ctx, filters)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func parseDTEFilters(r *http.Request) (*dte.DTEFilters, error) {
	filters := &dte.DTEFilters{
		IncludeAll: r.URL.Query().Get("all") == "true",
	}

	// 1. Si se incluyen todos los documentos no se establece el ID de la sucursal
	if !filters.IncludeAll {
		filters.BranchID = r.Context().Value("claims").(*models.AuthClaims).BranchID
	}

	// 2. Establecer los filtros de la request a la estructura de filtros
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	// 2.1 Fechas
	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, shared_error.NewGeneralServiceError("ListDTEsUseCase", "parseDTEFilters", "Invalid query param startDate format", nil)
		}
		startDate = &parsedStartDate
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, shared_error.NewGeneralServiceError("ListDTEsUseCase", "parseDTEFilters", "Invalid query param endDate format", nil)
		}
		endDate = &parsedEndDate
	}

	// 2.2 Status
	if status := r.URL.Query().Get("status"); status != "" {
		if !constants.ValidReceiverDocumentStates[strings.ToUpper(status)] {
			return nil, shared_error.NewGeneralServiceError("ListDTEsUseCase", "parseDTEFilters", "Invalid query param status, only 'received', 'invalidated' or 'rejected' are allowed", nil)
		} else {
			filters.Status = strings.ToUpper(status)
		}
	}

	// 2.3 Paginación
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		} else {
			filters.Page = 1 // Default
		}
	} else {
		filters.Page = 1 // Default
	}

	// 2.4 Transmisión
	if transmission := r.URL.Query().Get("transmission"); transmission != "" {
		if !constants.ValidTransmissionTypes[strings.ToUpper(transmission)] {
			return nil, shared_error.NewGeneralServiceError("ListDTEsUseCase", "parseDTEFilters", "Invalid query param transmission type, only 'normal' or 'contingency' are allowed", nil)
		} else {
			filters.Transmission = strings.ToUpper(transmission)
		}
	}

	// 2.5 Tamaño de página
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		} else {
			filters.PageSize = 5 // Default
		}
	} else {
		filters.PageSize = 5 // Default
	}

	// 2.6 Tipo de DTE
	if dteType := r.URL.Query().Get("type"); dteType != "" {
		if !constants.ValidDTETypes[dteType] {
			return nil, shared_error.NewGeneralServiceError("ListDTEsUseCase", "parseDTEFilters", "Invalid query param DTE type", nil)
		} else {
			filters.DTEType = dteType
		}
	}

	filters.StartDate = startDate
	filters.EndDate = endDate

	return filters, nil
}
