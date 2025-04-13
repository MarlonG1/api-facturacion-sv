package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	_ "github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type DTEHandler struct {
	invoiceUseCase      *dte.InvoiceUseCase
	ccfUseCase          *dte.CCFUseCase
	retentionUseCase    *dte.RetentionUseCase
	dteConsultUseCase   *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	contingencyHandler  *helpers.ContingencyHandler
	respWriter          *response.ResponseWriter
}

func NewDTEHandler(invoiceUseCase *dte.InvoiceUseCase, ccfUseCase *dte.CCFUseCase, retentionUseCase *dte.RetentionUseCase,
	dteConsultUseCase *dte.DTEConsultUseCase, invalidationUseCase *dte.InvalidationUseCase,
	contingencyHandler *helpers.ContingencyHandler) *DTEHandler {
	return &DTEHandler{
		invoiceUseCase:      invoiceUseCase,
		ccfUseCase:          ccfUseCase,
		retentionUseCase:    retentionUseCase,
		dteConsultUseCase:   dteConsultUseCase,
		invalidationUseCase: invalidationUseCase,
		contingencyHandler:  contingencyHandler,
		respWriter:          response.NewResponseWriter(),
	}
}

// CreateInvoice godoc
// @Summary      Create Invoice
// @Description  Create a new invoice
// @Tags         DTE
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param invoice body structs.CreateInvoiceRequest true "Invoice data"
// @Success      201 {object} response.APIDTEResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte/invoices [post]
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
		logs.Warn("Error transmitting invoice because", map[string]interface{}{"error": err.Error()})
		err = h.handleErrorForContingency(r.Context(), resp, constants.FacturaElectronica, options, err, w)
		if err != nil {
			h.respWriter.HandleError(w, err)
			return
		}
		return
	}

	// 3. Responder con la respuesta de la creación de invoice
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

// CreateCCF godoc
// @Summary      Create CCF
// @Description  Create a new CCF
// @Tags         DTE
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param ccf body structs.CreateCreditFiscalRequest true "CCF data"
// @Success      201 {object} response.APIDTEResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte/ccf [post]
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
		logs.Warn("Error transmitting CCF because", map[string]interface{}{"error": err.Error()})
		// 2.1. Verificar si aplica a contingencia
		err = h.handleErrorForContingency(r.Context(), resp, constants.CCFElectronico, options, err, w)
		if err != nil {
			h.respWriter.HandleError(w, err)
			return
		}
		return
	}

	// 3. Responder con la respuesta de la creación de CCF
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

// CreateRetention godoc
// @Summary      Create withholding certificate
// @Description  Create a new withholding certificate
// @Tags         DTE
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param retention body structs.CreateRetentionRequest true "Withholding certificate data"
// @Success      201 {object} response.APIDTEResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte/retention [post]
func (h *DTEHandler) CreateRetention(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de creación de retención a un DTO de solicitud
	var req structs.CreateRetentionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de creación de retención
	resp, options, err := h.retentionUseCase.Create(r.Context(), &req)
	if err != nil {
		logs.Warn("Error transmitting retention because", map[string]interface{}{"error": err.Error()})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder con la respuesta de la creación de retención
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

// GetByGenerationCode godoc
// @Summary      Get DTE by Generation Code
// @Description  Get DTE by Generation Code
// @Tags         DTE
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param id path string true "Generation Code"
// @Success 200 {object} dte.DTEResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte/{id} [get]
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

// GetAll godoc
// @Summary      Get All DTEs
// @Description  Get all DTEs
// @Tags         DTE
// @Accept       json
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param all query bool false "Include all DTEs"
// @Param startDate query string false "Start date in RFC3339 format"
// @Param endDate query string false "End date in RFC3339 format"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "DTE status"
// @Param transmission query string false "Transmission status"
// @Param type query string false "DTE type"
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dte.DTEListResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte [get]
func (h *DTEHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todos los DTEs ejecutando el caso de uso
	dtes, err := h.dteConsultUseCase.GetAllDTEs(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, dtes, nil)
}

// InvalidateDocument godoc
// @Summary      Invalidate Document
// @Description  Invalidate a document
// @Tags         DTE
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param request body structs.InvalidationRequest true "Invalidation request"
// @Success      200 {object} response.APIResponse
// @Failure      400 {object} response.APIResponse
// @Failure      500 {object} response.APIError
// @Router       /api/v1/dte/invalidate [post]
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

func (h *DTEHandler) handleErrorForContingency(ctx context.Context, dte interface{}, dteType string, options *response.SuccessOptions, err error, w http.ResponseWriter) error {
	// 1. Verificar si aplica a contingencia
	logs.Warn("Error transmitting DTE because", map[string]interface{}{
		"error": err.Error(),
	})

	contiType, reason := h.contingencyHandler.HandleContingency(ctx, dte, dteType, err)
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
