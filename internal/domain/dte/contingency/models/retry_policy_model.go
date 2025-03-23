package models

import "time"

type RetryPolicy struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	BackoffFactor   float64
}
