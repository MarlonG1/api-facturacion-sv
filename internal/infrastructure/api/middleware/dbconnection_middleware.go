package middleware

import (
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type DBConnectionMiddleware struct {
	connection *drivers.DbConnection
}

func NewDBConnectionMiddleware(connection *drivers.DbConnection) *DBConnectionMiddleware {
	return &DBConnectionMiddleware{connection: connection}
}

func (m *DBConnectionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := m.connection.Db.DB()

		if err != nil {
			logs.Error("Error connecting to database", map[string]interface{}{
				"error": err.Error(),
			})
			http.Error(w, "Error connecting to database", http.StatusInternalServerError)
			return
		}
		sqlDB.Ping()

		logs.Debug("Database connection established")
		next.ServeHTTP(w, r)
	})
}
