package dte

import "time"

// DTEDocument representa la relación entre un DTE y una sucursal
type DTEDocument struct {
	ID         string    `json:"id,omitempty"`
	BranchID   uint      `json:"branch_id"`
	DocumentID string    `json:"document_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`

	Details *DTEDetails `json:"dte_details,omitempty"`
}
