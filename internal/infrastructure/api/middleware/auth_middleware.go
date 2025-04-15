package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type AuthMiddleware struct {
	tokenService ports.TokenManager
	respWriter   *response.ResponseWriter
}

// NewAuthMiddleware crea una nueva instancia de AuthMiddleware. Recibe un servicio de tokens.
func NewAuthMiddleware(tokenService ports.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		respWriter:   response.NewResponseWriter(),
	}
}

// Handle es un middleware que valida el token de autorización.
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logs.Warn("Missing Authorization header", map[string]interface{}{
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logs.Warn("Invalid Authorization header format", map[string]interface{}{
				"header": authHeader,
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, "Invalid authorization format", nil)
			return
		}

		claims, err := m.tokenService.ValidateToken(parts[1])
		if err != nil {
			logs.Warn("Invalid token", map[string]interface{}{
				"error":  err.Error(),
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, "Invalid token", []string{"Token is expired or has been tampered with"})
			return
		}

		logs.Info("Token validated successfully", map[string]interface{}{
			"userID": claims.ClientID,
			"path":   r.URL.Path,
			"method": r.Method,
		})

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
