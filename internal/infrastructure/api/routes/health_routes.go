package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

// RegisterHealthRoutes registra las rutas de salud en el router
func RegisterHealthRoutes(router *mux.Router, healthHandler *handlers.HealthHandler) {
	router.HandleFunc("/health", healthHandler.CheckHealth).Methods("GET")
}
