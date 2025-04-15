package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap/containers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/routes"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	container *containers.Container
	srv       *http.Server

	privatePath string
	publicPath  string
}

// Initialize crea una nueva instancia de Server
func Initialize(container *containers.Container) *Server {
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
	s.configureProtectedMiddlewares(protected)

	s.router.Use(s.container.Middleware().DBConnectionMiddleware().Handler)
	s.configurePublicRoutes(public)
	s.configureProtectedRoutes(protected)

	logs.Info("Routes configured successfully", map[string]interface{}{
		"publicPath":    "/api/v1",
		"protectedPath": "/api/v1",
	})
}

func (s *Server) configureProtectedRoutes(protected *mux.Router) {
	routes.RegisterDTERoutes(protected, s.container.Handlers().DTEHandler())
	routes.RegisterMetricsRoutes(protected, s.container.Handlers().MetricsHandler())
}

func (s *Server) configureGlobalOptions() {
	s.router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) configurePublicRoutes(public *mux.Router) {
	routes.RegisterPublicAuthRoutes(public, s.container.Handlers().AuthHandler())
	routes.RegisterHealthRoutes(public, s.container.Handlers().HealthHandler())
	routes.RegisterTestRoutes(public, s.container.Handlers().TestHandler())
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
	s.ConfigureRoutes()

	s.srv = &http.Server{
		Handler:      s.router,
		Addr:         ":" + config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logs.Info("Server starting", map[string]interface{}{
		"port": config.Server.Port,
	})

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.Error("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logs.Info("Shutting down server...", nil)

	if s.srv == nil {
		logs.Warn("Server was not started, nothing to shutdown", nil)
		return nil
	}

	// Shutdown espera a que se completen las conexiones existentes
	return s.srv.Shutdown(ctx)
}
