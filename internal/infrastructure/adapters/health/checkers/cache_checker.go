package checkers

import (
	"context"
	"github.com/go-redis/redis/v8"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type redisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client *redis.Client) ports.ComponentChecker {
	return &redisChecker{client: client}
}

func (c *redisChecker) Name() string {
	return "cache"
}

func (c *redisChecker) Check() models.Health {
	ctx := context.Background()
	if err := c.client.Ping(ctx).Err(); err != nil {
		logs.Error("Cache connection error", map[string]interface{}{
			"error": err.Error(),
		})
		return models.Health{
			Status:  constants.StatusDown,
			Details: "Cache service unavailable",
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: "Cache service available",
	}
}
