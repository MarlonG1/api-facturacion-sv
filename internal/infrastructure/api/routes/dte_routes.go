package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"

	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDTERoutes(r *mux.Router, h *handlers.DTEHandler) {
	// Rutas para manejo de DTE
	r.HandleFunc("/dte/invoices", h.CreateInvoice).Methods(http.MethodPost)
	r.HandleFunc("/dte/ccf", h.CreateCCF).Methods(http.MethodPost)
	r.HandleFunc("/dte/invalidation", h.InvalidateDocument).Methods(http.MethodPost)
	r.HandleFunc("/dte/retention", h.CreateRetention).Methods(http.MethodPost)
	r.HandleFunc("/dte/credit_note", h.CreateCreditNote).Methods(http.MethodPost)

	// Rutas de consulta de DTE
	r.HandleFunc("/dte/{id}", h.GetByGenerationCode).Methods(http.MethodGet)
	r.HandleFunc("/dte", h.GetAll).Methods(http.MethodGet)
}
