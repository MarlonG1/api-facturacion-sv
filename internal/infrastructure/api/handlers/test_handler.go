package handlers

import (
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type TestHandler struct {
	testManager    test_endpoint.TestManager
	responseWriter *response.ResponseWriter
}

func NewTestHandler(testManager test_endpoint.TestManager) *TestHandler {
	return &TestHandler{
		testManager:    testManager,
		responseWriter: response.NewResponseWriter(),
	}
}

// RunSystemTest godoc
// @Summary      Run system test
// @Description  Run system test
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200 {object} models.TestResult
// @Failure      500 {object} response.APIError
// @Router       /api/v1/test [get]
func (h *TestHandler) RunSystemTest(w http.ResponseWriter, r *http.Request) {
	result, err := h.testManager.RunSystemTest()
	if err != nil {
		logs.Error("System test failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.responseWriter.Error(w, http.StatusInternalServerError, "System test failed", nil)
		return
	}

	h.responseWriter.Success(w, http.StatusOK, result, nil)
}
