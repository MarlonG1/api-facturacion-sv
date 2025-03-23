package models

import (
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/models"
	"time"
)

// BatchConfig configuración específica para el servicio de contingencia
type BatchConfig struct {
	Ambient       string
	BatchSize     int
	RetryInterval time.Duration
	MaxInterval   time.Duration
	BackoffFactor float64
}

// GetAmbient obtiene el ambiente configurado
func (c *BatchConfig) GetAmbient() string {
	return c.Ambient
}

// GetBatchSize obtiene el tamaño de lote configurado
func (c *BatchConfig) GetBatchSize() int {
	return c.BatchSize
}

// GetRetryInterval obtiene el intervalo de reintento inicial
func (c *BatchConfig) GetRetryInterval() time.Duration {
	return c.RetryInterval
}

// GetMaxInterval obtiene el intervalo máximo de reintento
func (c *BatchConfig) GetMaxInterval() time.Duration {
	return c.MaxInterval
}

// GetBackoffFactor obtiene el factor de crecimiento para backoff exponencial
func (c *BatchConfig) GetBackoffFactor() float64 {
	return c.BackoffFactor
}

// GetRetryPolicy construye una política de reintentos
func (c *BatchConfig) GetRetryPolicy() models.RetryPolicy {
	return models.RetryPolicy{
		MaxAttempts:     3, // Valor por defecto
		InitialInterval: c.RetryInterval,
		MaxInterval:     c.MaxInterval,
		BackoffFactor:   c.BackoffFactor,
	}
}

// NewContingencyConfig constructor para BatchConfig
func NewContingencyConfig(retryInterval, maxInterval time.Duration, backoffFactor float64) *BatchConfig {
	return &BatchConfig{
		Ambient:       env.Server.AmbientCode,
		BatchSize:     env.Server.MaxBatchSize,
		RetryInterval: retryInterval,
		MaxInterval:   maxInterval,
		BackoffFactor: backoffFactor,
	}
}
