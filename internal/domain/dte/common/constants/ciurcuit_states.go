package constants

type State int32

const (
	StateClosed   State = iota // Representa el estado cerrado del circuit breaker
	StateOpen                  // Representa el estado abierto del circuit breaker
	StateHalfOpen              // Representa el estado semi-abierto del circuit breaker
)
