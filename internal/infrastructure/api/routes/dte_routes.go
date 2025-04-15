package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"

	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDTERoutes(r *mux.Router, h *handlers.DTEHandler) {
	// Rutas para manejo de DTE
	for path, _ := range h.GenericHandler.GetDocumentConfigs() {
		r.HandleFunc(path, h.GenericHandler.HandleCreate).Methods(http.MethodPost)
	}

	// Rutas de consulta de DTE e Invalidación
	r.HandleFunc("/dte/invalidation", h.InvalidateDocument).Methods(http.MethodPost)
	r.HandleFunc("/dte/{id}", h.GetByGenerationCode).Methods(http.MethodGet)
	r.HandleFunc("/dte", h.GetAll).Methods(http.MethodGet)
}
