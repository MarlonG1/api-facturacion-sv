package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

func RegisterTestRoutes(router *mux.Router, testHandler *handlers.TestHandler) {
	router.HandleFunc("/test", testHandler.RunSystemTest).Methods("GET")
}
