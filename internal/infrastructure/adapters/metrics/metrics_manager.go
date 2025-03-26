package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/models"
	metricsPort "github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type MetricManager struct {
	cache     ports.CacheManager
	endpoints []struct {
		path   string
		method string
	}
}

func NewMetricManager(cache ports.CacheManager) metricsPort.MetricsManager {
	return &MetricManager{
		cache: cache,
		endpoints: []struct {
			path   string
			method string
		}{
			{path: "invoices", method: "POST"},
			{path: "ccf", method: "POST"},
			{path: "invalidation", method: "POST"},
			{path: "dte", method: "GET"},
			{path: "dte/{id}", method: "GET"},
		},
	}
}

func (m *MetricManager) GetEndpointMetrics(systemNIT, method, endpoint string) (*models.EndpointMetrics, error) {
	durationsKey := fmt.Sprintf("metrics:%s:%s:%s:durations", systemNIT, method, endpoint)
	countersKey := fmt.Sprintf("metrics:%s:%s:%s:counters", systemNIT, method, endpoint)

	// Obtener histórico de duraciones
	results, err := m.cache.LRange(durationsKey, 0, 19)
	if err != nil {
		return nil, err
	}

	// Obtener contadores globales
	var counters models.EndpointCounters
	countersData, err := m.cache.Get(countersKey)
	if err == nil {
		json.Unmarshal([]byte(countersData), &counters)
	}

	metrics := &models.EndpointMetrics{
		Path:          endpoint,
		Method:        method,
		SystemNIT:     systemNIT,
		TotalRequests: counters.TotalRequests,
		SuccessCount:  counters.SuccessCount,
		ErrorCount:    counters.ErrorCount,
	}

	// Procesar histórico de duraciones
	var totalDuration int64
	minDuration := int64(^uint64(0) >> 1) // Max int64
	maxDuration := int64(0)

	for _, result := range results {
		var metric models.RequestMetric
		if err := json.Unmarshal([]byte(result), &metric); err != nil {
			continue
		}

		metrics.LastDurations = append(metrics.LastDurations, metric.Duration)
		totalDuration += metric.Duration

		if metric.Duration < minDuration {
			minDuration = metric.Duration
		}
		if metric.Duration > maxDuration {
			maxDuration = metric.Duration
		}
	}

	// Calcular promedio del histórico reciente
	if len(metrics.LastDurations) > 0 {
		metrics.CurrentAverage = totalDuration / int64(len(metrics.LastDurations))
		metrics.MinDuration = minDuration
		metrics.MaxDuration = maxDuration
	}

	return metrics, nil
}

func (m *MetricManager) GetAllMetricsEndpoint(systemNIT string) (map[string]*models.EndpointMetrics, error) {
	// Lista de endpoints conocidos del sistema
	allMetrics := make(map[string]*models.EndpointMetrics)

	for _, ep := range m.endpoints {
		metrics, err := m.GetEndpointMetrics(systemNIT, ep.method, ep.path)
		if err != nil {
			logs.Warn("Failed to get metrics for endpoint", map[string]interface{}{
				"endpoint": ep.path,
				"method":   ep.method,
				"error":    err.Error(),
			})
			continue
		}

		key := fmt.Sprintf("%s-%s", ep.method, ep.path)
		allMetrics[key] = metrics
	}

	return allMetrics, nil
}
