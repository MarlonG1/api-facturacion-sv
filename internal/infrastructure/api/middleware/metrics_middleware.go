package middleware

import (
	"encoding/json"
	"fmt"
	metricsModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/metrics/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"net/http"
	"strings"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type MetricsMiddleware struct {
	cache      ports.CacheManager
	maxMetrics int
}

func NewMetricsMiddleware(cache ports.CacheManager) *MetricsMiddleware {
	return &MetricsMiddleware{
		cache:      cache,
		maxMetrics: 20,
	}
}

func extractEndpoint(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return path
}

func (m *MetricsMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
			written:        false,
		}

		systemNIT := r.Context().Value("claims").(*models.AuthClaims).NIT

		start := utils.TimeNow()
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		endpoint := extractEndpoint(r.URL.Path)
		durationsKey := fmt.Sprintf("metrics:%s:%s:%s:durations", systemNIT, r.Method, endpoint)
		countersKey := fmt.Sprintf("metrics:%s:%s:%s:counters", systemNIT, r.Method, endpoint)

		metric := metricsModels.RequestMetric{
			Duration:   duration.Milliseconds(),
			Timestamp:  utils.TimeNow(),
			StatusCode: rw.status,
		}

		metricData, _ := json.Marshal(metric)

		if err := m.cache.LPush(durationsKey, metricData); err != nil {
			logs.Error("Failed to store metric", map[string]interface{}{
				"error": err.Error(),
			})
		}

		if err := m.cache.LTrim(durationsKey, 0, 19); err != nil {
			logs.Error("Failed to trim metrics", map[string]interface{}{
				"error": err.Error(),
			})
		}

		var counters metricsModels.EndpointCounters
		if countersData, err := m.cache.Get(countersKey); err == nil {
			json.Unmarshal([]byte(countersData), &counters)
		}

		counters.TotalRequests++
		if rw.status >= 200 && rw.status < 300 {
			counters.SuccessCount++
		} else {
			counters.ErrorCount++
		}

		// Guardar contadores actualizados
		if countersData, err := json.Marshal(counters); err == nil {
			if err := m.cache.Set(countersKey, countersData, 24*time.Hour); err != nil {
				logs.Error("Failed to update counters", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.status = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.status = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}
