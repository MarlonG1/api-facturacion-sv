package metrics

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/models"
)

// MetricsManager es una interfaz que define los métodos para obtener y registrar métricas
type MetricsManager interface {
	// GetAllMetricsEndpoint obtiene las métricas actuales
	GetAllMetricsEndpoint(systemNIT string) (map[string]*models.EndpointMetrics, error)
	// GetEndpointMetrics obtiene las métricas de un endpoint específico
	GetEndpointMetrics(systemNIT, method, endpoint string) (*models.EndpointMetrics, error)
}
