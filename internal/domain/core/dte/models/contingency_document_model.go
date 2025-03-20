package models

import "time"

// ContingencyDocument representa un documento en estado de contingencia
type ContingencyDocument struct {
	ID              string    `json:"id,omitempty"`
	DocumentID      string    `json:"document_id"`
	BranchID        uint      `json:"branch_id"`
	ContingencyType int8      `json:"contingency_type"`
	Reason          string    `json:"reason"`
	BatchID         string    `json:"batch_id"`
	MHBatchID       string    `json:"mh_batch_id,omitempty"`
	Observations    string    `json:"observations,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
