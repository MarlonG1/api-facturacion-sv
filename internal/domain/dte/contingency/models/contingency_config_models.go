package models

import "time"

type ContingencyResult struct {
	ContingencyType   int
	ContingencyReason string
	ShouldRetry       bool
	RetryConfig       *RetryConfig
}

type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
}

type ErrClassification struct {
	Type       int
	Reason     string
	Retryable  bool
	RetryDelay time.Duration
}
