package batch

import (
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"time"
)

// ProductionConfig es la configuración de contingencia para producción
type ProductionConfig struct{}

func (c *ProductionConfig) GetMaxRetries() int              { return 3 }
func (c *ProductionConfig) GetRetryInterval() time.Duration { return 5 * time.Second }
func (c *ProductionConfig) GetMaxInterval() time.Duration   { return 2 * time.Minute }
func (c *ProductionConfig) GetBatchSize() int               { return env.Server.MaxBatchSize }
func (c *ProductionConfig) GetAmbient() string              { return env.Server.AmbientCode }
func (c *ProductionConfig) GetBackoffFactor() float64       { return 2.0 }

// TestConfig es la configuración de contingencia para pruebas
type TestConfig struct{}

func (c *TestConfig) GetMaxRetries() int              { return 1 }
func (c *TestConfig) GetRetryInterval() time.Duration { return 1 * time.Millisecond }
func (c *TestConfig) GetMaxInterval() time.Duration   { return 2 * time.Millisecond }
func (c *TestConfig) GetBatchSize() int               { return 1 }
func (c *TestConfig) GetAmbient() string              { return "00" }
func (c *TestConfig) GetBackoffFactor() float64       { return 1.0 }

// RealTimeProvider es un proveedor de tiempo real
type RealTimeProvider struct{}

func (p *RealTimeProvider) Now() time.Time {
	return utils.TimeNow()
}

func (p *RealTimeProvider) Sleep(d time.Duration) {
	time.Sleep(d)
}
