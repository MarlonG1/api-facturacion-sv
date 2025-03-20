package models

type HealthStatus struct {
	Status     string            `json:"status"`
	Components map[string]Health `json:"components"`
	Timestamp  string            `json:"timestamp"`
}

type Health struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}
