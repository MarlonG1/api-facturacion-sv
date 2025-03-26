package handlers

import (
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"net/http"
)

type DTEHandler struct {
	invoiceUseCase      *dte.InvoiceUseCase
	ccfUseCase          *dte.CCFUseCase
	dteConsultUseCase   *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	contingencyHandler  *helpers.ContingencyHandler
	respWriter          *response.ResponseWriter
}

func NewDTEHandler(invoiceUseCase *dte.InvoiceUseCase, ccfUseCase *dte.CCFUseCase, dteConsultUseCase *dte.DTEConsultUseCase, invalidationUseCase *dte.InvalidationUseCase, contingencyHandler *helpers.ContingencyHandler) *DTEHandler {
	return &DTEHandler{
		invoiceUseCase:      invoiceUseCase,
		ccfUseCase:          ccfUseCase,
		dteConsultUseCase:   dteConsultUseCase,
		invalidationUseCase: invalidationUseCase,
		contingencyHandler:  contingencyHandler,
		respWriter:          response.NewResponseWriter(),
	}
}

func (h *DTEHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de creación de invoice a un DTO de solicitud
	var req structs.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de creación de factura electrónica
	resp, options, err := h.invoiceUseCase.Create(r.Context(), &req)
	if err != nil {
		// 2.1. Verificar si aplica a contingencia
		logs.Warn("Error transmitting CCF because", map[string]interface{}{"error": err.Error()})
		contiType, reason := h.contingencyHandler.HandleContingency(r.Context(), resp, constants.CCFElectronico, err)

		if contiType == nil || reason == nil {
			logs.Error("Error creating CCF", map[string]interface{}{"error": err.Error()})
			h.respWriter.HandleError(w, err)
			return
		}

		// 2.2. Actualizar la identificación de contingencia en el JSON del DTE
		err = h.handleErrorForContingency(r.Context(), resp, options, err, w)
		if err != nil {
			h.respWriter.HandleError(w, err)
			return
		}
		return
	}

	// 3. Responder con la respuesta de la creación de invoice
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

func (h *DTEHandler) CreateCCF(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de creación de CCF a un DTO de solicitud
	var req structs.CreateCreditFiscalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de creación de CCF
	resp, options, err := h.ccfUseCase.Create(r.Context(), &req)
	if err != nil {
		// 2.1. Verificar si aplica a contingencia
		logs.Warn("Error transmitting CCF because", map[string]interface{}{"error": err.Error()})
		contiType, reason := h.contingencyHandler.HandleContingency(r.Context(), resp, constants.CCFElectronico, err)

		if contiType == nil || reason == nil {
			logs.Error("Error creating CCF", map[string]interface{}{"error": err.Error()})
			h.respWriter.HandleError(w, err)
			return
		}

		// 2.2. Actualizar la identificación de contingencia en el JSON del DTE
		err = h.handleErrorForContingency(r.Context(), resp, options, err, w)
		if err != nil {
			h.respWriter.HandleError(w, err)
			return
		}

		return
	}

	// 3. Responder con la respuesta de la creación de CCF
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

func (h *DTEHandler) GetDTEByGenerationCode(w http.ResponseWriter, r *http.Request) {
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

func (h *DTEHandler) InvalidateDocument(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de invalidación de documento a un DTO de solicitud
	var req structs.InvalidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de invalidación de documento
	err := h.invalidationUseCase.InvalidateDocument(r.Context(), req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "DTE invalidated successfully", nil)
}

func (h *DTEHandler) handleErrorForContingency(ctx context.Context, dte interface{}, options *response.SuccessOptions, err error, w http.ResponseWriter) error {
	// 1. Verificar si aplica a contingencia
	logs.Warn("Error transmitting DTE because", map[string]interface{}{
		"error": err.Error(),
	})

	contiType, reason := h.contingencyHandler.HandleContingency(ctx, dte, constants.CCFElectronico, err)
	if contiType == nil || reason == nil {
		logs.Error("Error creating DTE contingency", map[string]interface{}{"error": err.Error()})
		return err
	}

	// 2. Actualizar la identificación de contingencia en el JSON del DTE
	updatedDTE, err := utils.UpdateContingencyIdentification(dte, contiType, reason)
	if err != nil {
		return err
	}

	// 3. Responder con la respuesta de la creación del DTE
	h.respWriter.Success(w, http.StatusCreated, updatedDTE, options)
	return nil
}
