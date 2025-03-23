package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

func RegisterCCFRoutes(r *mux.Router, handler *handlers.CCFHandler) {
	r.HandleFunc("/ccf", handler.CreateCCF).Methods("POST")
}
