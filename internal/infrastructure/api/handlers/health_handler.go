package handlers

import (
	_ "github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"net/http"
)

type HealthHandler struct {
	healthManager  ports.HealthManager
	responseWriter *response.ResponseWriter
}

func NewHealthHandler(checkHealthUseCase ports.HealthManager) *HealthHandler {
	return &HealthHandler{
		healthManager:  checkHealthUseCase,
		responseWriter: response.NewResponseWriter(),
	}
}

// CheckHealth godoc
// @Summary      Health Check
// @Description  Check the health of all core service
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} models.HealthStatus
// @Failure      500 {object} response.APIError
// @Router       /api/v1/health [get]
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	logs.Info("Starting health check")
	defer logs.Info("Health check finished")

	status, err := h.healthManager.CheckHealth()
	if err != nil {
		logs.Error("Health check failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.responseWriter.Error(w, http.StatusInternalServerError, "Health check failed", []string{err.Error()})
		return
	}

	h.responseWriter.Success(w, http.StatusOK, status, nil)
}
