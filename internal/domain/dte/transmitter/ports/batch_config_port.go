package ports

import "time"

// BatchConfiguration es una interfaz que define los métodos que debe implementar una configuración de contingencia
type BatchConfiguration interface {
	GetMaxRetries() int              // GetMaxRetries retorna el número máximo de reintentos
	GetRetryInterval() time.Duration // GetRetryInterval retorna el intervalo de reintentos
	GetMaxInterval() time.Duration   // GetMaxInterval retorna el intervalo máximo de reintentos
	GetBatchSize() int               // GetBatchSize retorna el tamaño del lote
	GetAmbient() string              // GetAmbient retorna el ambiente
	GetBackoffFactor() float64       // GetBackoffFactor retorna el factor de retroceso
}
