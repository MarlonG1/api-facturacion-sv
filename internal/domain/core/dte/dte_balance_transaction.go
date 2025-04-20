package dte

import "time"

type BalanceTransaction struct {
	BalanceControlID     uint      `json:"balance_control_id"`
	AdjustmentDocumentID string    `json:"adjustment_document_id"`
	TransactionType      string    `json:"transaction_type"`
	TaxedAmount          float64   `json:"taxed_amount"`
	ExemptAmount         float64   `json:"exempt_amount"`
	NotSubjectAmount     float64   `json:"non_subject_amount"`
	CreatedAt            time.Time `json:"created_at"`

	AdjustmentDocument *DTEDetails     `json:"adjustment_document,omitempty"`
	BalanceControl     *BalanceControl `json:"balance_control,omitempty"`
}
