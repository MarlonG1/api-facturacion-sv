package middleware

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	responseWriter *response.ResponseWriter
}

func NewTimeoutMiddleware() *TimeoutMiddleware {
	return &TimeoutMiddleware{
		responseWriter: response.NewResponseWriter(),
	}
}

func (m *TimeoutMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 14*time.Second)
		defer cancel()

		r = r.WithContext(ctx)
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Canal para esperar la finalización de la solicitud
		done := make(chan struct{})
		go func() {
			next.ServeHTTP(rw, r.WithContext(ctx))
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			timeoutTitle := i18n.TranslateServiceArgs("RequestTimeOutTitle")
			timeoutMessage := i18n.TranslateServiceArgs("RequestTimeOut")

			logs.Warn("Request timed out", map[string]interface{}{
				"method":              r.Method,
				"path":                r.URL.Path,
				"error":               ctx.Err(),
				"error shown":         timeoutMessage,
				"error message shown": timeoutTitle,
			})
			if !rw.written {
				m.responseWriter.Error(
					w,
					http.StatusRequestTimeout,
					timeoutTitle,
					[]string{timeoutMessage},
				)
			}
		}
	})
}
