package handlers

import (
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type MetricsHandler struct {
	metricsManager metrics.MetricsManager
	responseWriter *response.ResponseWriter
}

func NewMetricsHandler(metricsManager metrics.MetricsManager) *MetricsHandler {
	return &MetricsHandler{
		metricsManager: metricsManager,
		responseWriter: response.NewResponseWriter(),
	}
}

// GetEndpointMetrics godoc
// @Summary      Get endpoint metrics
// @Description  Get metrics for a specific endpoint
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param Authorization header string true "Token JWT with Format 'Bearer {token}'"
// @Param endpoint query string false "Endpoint to filter metrics"
// @Param method query string false "HTTP method to filter metrics"
// @Success      200 {object} models.EndpointMetrics
// @Failure      400 {object} response.APIError
// @Failure      500 {object} response.APIError
// @Router       /api/v1/metrics [get]
func (h *MetricsHandler) GetEndpointMetrics(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.AuthClaims)

	endpoint := r.URL.Query().Get("endpoint")
	method := r.URL.Query().Get("method")

	// Si no hay filtros, obtenemos todas las métricas de endpoints
	if endpoint == "" || method == "" {
		allMetrics, err := h.metricsManager.GetAllMetricsEndpoint(claims.NIT)
		if err != nil {
			logs.Error("Failed to get all endpoint endpointMetrics", map[string]interface{}{
				"error":     err.Error(),
				"systemNIT": claims.NIT,
			})
			h.responseWriter.Error(w, http.StatusInternalServerError, "Failed to get endpoint endpointMetrics", nil)
			return
		}
		h.responseWriter.Success(w, http.StatusOK, allMetrics, nil)
		return
	}

	// Si hay filtros, obtenemos las métricas específicas
	endpointMetrics, err := h.metricsManager.GetEndpointMetrics(claims.NIT, method, endpoint)
	if err != nil {
		logs.Error("Failed to get endpoint endpointMetrics", map[string]interface{}{
			"error":     err.Error(),
			"systemNIT": claims.NIT,
			"endpoint":  endpoint,
			"method":    method,
		})
		h.responseWriter.Error(w, http.StatusInternalServerError, "Failed to get endpoint endpointMetrics", nil)
		return
	}

	h.responseWriter.Success(w, http.StatusOK, endpointMetrics, nil)
}
