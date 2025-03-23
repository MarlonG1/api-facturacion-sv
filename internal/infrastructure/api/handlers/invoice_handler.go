package handlers

import (
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"net/http"
)

type InvoiceHandler struct {
	invoiceUseCase *dte.InvoiceUseCase

	respWriter *response.ResponseWriter
}

func NewInvoiceHandler(invoiceUseCase *dte.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{
		respWriter:     response.NewResponseWriter(),
		invoiceUseCase: invoiceUseCase,
	}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de creación de invoice a un DTO de solicitud
	var req structs.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de creación de invoice
	resp, _, err := h.invoiceUseCase.Create(r.Context(), &req)
	if err != nil {
		logs.Error("Error creating invoice", map[string]interface{}{"error": err.Error()})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder con la respuesta de la creación de invoice
	h.respWriter.Success(w, http.StatusCreated, resp, nil)
}
