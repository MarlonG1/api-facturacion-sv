package bootstrap

import "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/middleware"

type MiddlewareContainer struct {
	services *ServicesContainer

	corsMid   *middleware.CorsMiddleware
	authMid   *middleware.AuthMiddleware
	tokenMid  *middleware.TokenExtractor
	errorMid  *middleware.ErrorMiddleware
	metricMid *middleware.MetricsMiddleware
}

func NewMiddlewareContainer(services *ServicesContainer) *MiddlewareContainer {
	return &MiddlewareContainer{
		services: services,
	}
}

func (c *MiddlewareContainer) Initialize() {
	c.corsMid = middleware.NewCorsMiddleware(
		[]string{"*"},
		nil,
		nil,
	)
	c.tokenMid = middleware.NewTokenExtractor()
	c.errorMid = middleware.NewErrorMiddleware()
	c.authMid = middleware.NewAuthMiddleware(c.services.TokenManager())
	c.metricMid = middleware.NewMetricsMiddleware(c.services.CacheManager())
}

func (c *MiddlewareContainer) CorsMiddleware() *middleware.CorsMiddleware {
	return c.corsMid
}

func (c *MiddlewareContainer) AuthMiddleware() *middleware.AuthMiddleware {
	return c.authMid
}

func (c *MiddlewareContainer) TokenExtractor() *middleware.TokenExtractor {
	return c.tokenMid
}

func (c *MiddlewareContainer) ErrorMiddleware() *middleware.ErrorMiddleware {
	return c.errorMid
}

func (c *MiddlewareContainer) MetricsMiddleware() *middleware.MetricsMiddleware {
	return c.metricMid
}
