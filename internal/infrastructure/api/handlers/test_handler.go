package handlers

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint/ports"
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type TestHandler struct {
	testManager    ports.TestManager
	responseWriter *response.ResponseWriter
}

func NewTestHandler(testManager ports.TestManager) *TestHandler {
	return &TestHandler{
		testManager:    testManager,
		responseWriter: response.NewResponseWriter(),
	}
}

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
