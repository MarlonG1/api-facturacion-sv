package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

func RegisterDTERoutes(r *mux.Router, h *handlers.DTEHandler) {
	r.HandleFunc("/dte/invoices", h.CreateInvoice).Methods("POST")
	r.HandleFunc("/dte/ccf", h.CreateCCF).Methods("POST")
	r.HandleFunc("/dte/invalidation", h.InvalidateDocument).Methods("POST")
	r.HandleFunc("/dte/{id}", h.GetByGenerationCode).Methods("GET")
	r.HandleFunc("/dte", h.GetAll).Methods("GET")
}
