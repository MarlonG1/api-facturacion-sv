package ports

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/models"
)

// MetricsManager es una interfaz que define los métodos para obtener y registrar métricas
type MetricsManager interface {
	GetMetrics(ctx context.Context) (*models.Metrics, error) // GetMetrics obtiene las métricas actuales
}
