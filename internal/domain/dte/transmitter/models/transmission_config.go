package models

import (
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/models"
	"time"
)

// TransmissionConfig configuración específica para el servicio de contingencia
type TransmissionConfig struct {
	Ambient       string
	BatchSize     int
	RetryInterval time.Duration
	MaxInterval   time.Duration
	BackoffFactor float64
}

// GetAmbient obtiene el ambiente configurado
func (c *TransmissionConfig) GetAmbient() string {
	return c.Ambient
}

// GetBatchSize obtiene el tamaño de lote configurado
func (c *TransmissionConfig) GetBatchSize() int {
	return c.BatchSize
}

// GetRetryInterval obtiene el intervalo de reintento inicial
func (c *TransmissionConfig) GetRetryInterval() time.Duration {
	return c.RetryInterval
}

// GetMaxInterval obtiene el intervalo máximo de reintento
func (c *TransmissionConfig) GetMaxInterval() time.Duration {
	return c.MaxInterval
}

// GetBackoffFactor obtiene el factor de crecimiento para backoff exponencial
func (c *TransmissionConfig) GetBackoffFactor() float64 {
	return c.BackoffFactor
}

// GetRetryPolicy construye una política de reintentos
func (c *TransmissionConfig) GetRetryPolicy() models.RetryPolicy {
	return models.RetryPolicy{
		MaxAttempts:     3, // Valor por defecto
		InitialInterval: c.RetryInterval,
		MaxInterval:     c.MaxInterval,
		BackoffFactor:   c.BackoffFactor,
	}
}

func NewTransmissionConfig(retryInterval, maxInterval time.Duration, backoffFactor float64) *TransmissionConfig {
	return &TransmissionConfig{
		Ambient:       config.Server.AmbientCode,
		BatchSize:     config.Server.MaxBatchSize,
		RetryInterval: retryInterval,
		MaxInterval:   maxInterval,
		BackoffFactor: backoffFactor,
	}
}
