package server

import (
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/routes"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	container *bootstrap.Container

	privatePath string
	publicPath  string
}

// Initialize crea una nueva instancia de Server
func Initialize(container *bootstrap.Container) *Server {
	return &Server{
		router:      mux.NewRouter(),
		container:   container,
		publicPath:  "/api/v1",
		privatePath: "/api/v1",
	}
}

func (s *Server) ConfigureRoutes() {
	s.configureGlobalMiddlewares()
	s.configureGlobalOptions()

	// Configurar rutas públicas y protegidas
	public := s.router.PathPrefix(s.publicPath).Subrouter()
	protected := s.router.PathPrefix(s.privatePath).Subrouter()

	s.configurePublicRoutes(public)
	s.configureProtectedMiddlewares(protected)
	s.configureProtectedRoutes(protected)

	logs.Info("Routes configured successfully", map[string]interface{}{
		"publicPath":    "/api/v1",
		"protectedPath": "/api/v1",
	})
}

func (s *Server) configureProtectedRoutes(protected *mux.Router) {

}

func (s *Server) configureGlobalOptions() {
	s.router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) configurePublicRoutes(public *mux.Router) {
	routes.RegisterPublicAuthRoutes(public, s.container.Handlers().AuthHandler())
}

func (s *Server) configureGlobalMiddlewares() {
	s.router.Use(s.container.Middleware().CorsMiddleware().Handler)
	s.router.Use(s.container.Middleware().ErrorMiddleware().Handler)
}

func (s *Server) configureProtectedMiddlewares(protected *mux.Router) {
	protected.Use(s.container.Middleware().AuthMiddleware().Handle)
	protected.Use(s.container.Middleware().TokenExtractor().ExtractToken)
	protected.Use(s.container.Middleware().MetricsMiddleware().Handle)
}

func (s *Server) Start() error {
	err := s.container.Initialize()
	if err != nil {
		return err
	}

	s.ConfigureRoutes()

	srv := &http.Server{
		Handler:      s.router,
		Addr:         ":" + env.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logs.Info("Server starting", map[string]interface{}{
		"port": env.Server.Port,
	})

	if err := srv.ListenAndServe(); err != nil {
		logs.Error("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
