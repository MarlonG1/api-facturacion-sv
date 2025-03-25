package circuit

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"sync"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type CircuitBreaker struct {
	failures    int32
	lastFailure time.Time
	threshold   int32
	resetTime   time.Duration
	state       constants.State
	mu          sync.RWMutex
}

func NewCircuitBreaker(threshold int32, resetTime time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		resetTime: resetTime,
		state:     constants.StateClosed,
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case constants.StateClosed:
		return true
	case constants.StateOpen:
		if time.Since(cb.lastFailure) > cb.resetTime {
			logs.Info("Circuit breaker entering half-open state", map[string]interface{}{
				"lastFailure": cb.lastFailure,
				"resetTime":   cb.resetTime,
			})
			cb.state = constants.StateHalfOpen
			return true
		}
		logs.Debug("Circuit breaker is open, blocking request", map[string]interface{}{
			"remainingTime": cb.resetTime - time.Since(cb.lastFailure),
		})
		return false
	case constants.StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	if cb.state != constants.StateClosed {
		logs.Info("Circuit breaker closing after success", map[string]interface{}{
			"previousState": cb.state,
		})
	}
	cb.state = constants.StateClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = utils.TimeNow()

	if cb.failures >= cb.threshold && cb.state != constants.StateOpen {
		logs.Warn("Circuit breaker opening due to failures", map[string]interface{}{
			"failures":  cb.failures,
			"threshold": cb.threshold,
		})
		cb.state = constants.StateOpen
	}
}

func (cb *CircuitBreaker) GetState() constants.State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) GetFailureCount() int32 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}
