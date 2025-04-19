package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type DTEHandler struct {
	GenericHandler      *GenericCreatorDTEHandler
	dteConsultUseCase   *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	respWriter          *response.ResponseWriter
}

func NewDTEHandler(
	dteConsultUseCase *dte.DTEConsultUseCase,
	invalidationUseCase *dte.InvalidationUseCase,
	genericHandler *GenericCreatorDTEHandler,
) *DTEHandler {
	return &DTEHandler{
		GenericHandler:      genericHandler,
		dteConsultUseCase:   dteConsultUseCase,
		invalidationUseCase: invalidationUseCase,
		respWriter:          response.NewResponseWriter(),
	}
}

// NOTA IMPORTANTE:
// Ahora la creacion de un DTE se maneja en el GenericCreatorDTEHandler y sus rutas en el router

// GetByGenerationCode maneja la solicitud HTTP para obtener un DTE por su código de generación
func (h *DTEHandler) GetByGenerationCode(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el código de generación
	generationCode := helpers.GetRequestVar(r, "id")

	// 2. Obtener DTE ejecutando el caso de uso
	dte, err := h.dteConsultUseCase.GetByGenerationCode(r.Context(), generationCode)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, dte, nil)
}

// GetAll maneja la solicitud HTTP para obtener todos los DTEs
func (h *DTEHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todos los DTEs ejecutando el caso de uso
	dtes, err := h.dteConsultUseCase.GetAllDTEs(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, dtes, nil)
}

// InvalidateDocument maneja la solicitud HTTP para invalidar un DTE
func (h *DTEHandler) InvalidateDocument(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de invalidación de documento a un DTO de solicitud
	var req structs.CreateInvalidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de invalidación de documento
	invalidation, err := h.invalidationUseCase.InvalidateDocument(r.Context(), req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, invalidation, nil)
}
