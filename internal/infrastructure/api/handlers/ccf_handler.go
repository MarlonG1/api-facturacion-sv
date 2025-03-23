package handlers

import (
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"net/http"
)

type CCFHandler struct {
	ccfUseCase *dte.CCFUseCase
	respWriter *response.ResponseWriter
}

func NewCCFHandler(ccfUseCase *dte.CCFUseCase) *CCFHandler {
	return &CCFHandler{
		ccfUseCase: ccfUseCase,
		respWriter: response.NewResponseWriter(),
	}
}

func (h *CCFHandler) CreateCCF(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de creación de CCF a un DTO de solicitud
	var req structs.CreateCreditFiscalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de creación de CCF
	resp, _, err := h.ccfUseCase.Create(r.Context(), &req)
	if err != nil {
		logs.Error("Error creating CCF", map[string]interface{}{"error": err.Error()})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder con la respuesta de la creación de CCF
	h.respWriter.Success(w, http.StatusCreated, resp, nil)
}
