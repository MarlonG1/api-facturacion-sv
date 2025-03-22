package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

func RegisterPublicAuthRoutes(r *mux.Router, h *handlers.AuthHandler) {
	r.HandleFunc("/register", h.Register).Methods("POST")
}
