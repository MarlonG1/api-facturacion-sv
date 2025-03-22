package db_models

import "time"

type FailedSequenceNumber struct {
	ID                  uint      `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	BranchID            uint      `gorm:"column:branch_id;type:uint;not null;index:idx_failed_seq"`
	DTEType             string    `gorm:"column:dte_type;type:varchar(2);not null;index:idx_failed_seq"`
	SequenceNumber      uint      `gorm:"column:sequence_number;type:uint;not null;index:idx_failed_seq"`
	Year                uint      `gorm:"column:year;type:uint;not null;index:idx_failed_seq"`
	FailureReason       string    `gorm:"column:failure_reason;type:text;not null"`
	ResponseCode        string    `gorm:"column:response_code;type:varchar(10)"`
	OriginalRequestData string    `gorm:"column:original_request_data;type:json;not null"`
	MHResponse          string    `gorm:"column:mh_response;type:text"`
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (FailedSequenceNumber) TableName() string {
	return "failed_sequence_numbers"
}
