package models

// DomainEvent representa un evento de dominio
type DomainEvent struct {
	ID         uint   `json:"id,omitempty"`
	UserID     uint   `json:"user_id"`
	BranchID   uint   `json:"branch_id"`
	EventType  string `json:"event_type"`
	Payload    string `json:"payload"`
	OccurredAt string `json:"occurred_at"`
}
