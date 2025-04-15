package db_models

import "time"

// DTEBalanceControl representa el control de saldo de un DTE (Documento Tributario Electrónico).
// Se utiliza para almacenar el saldo de un DTE en la base de datos y recuperarlo para su procesamiento de una manera más eficiente
// lo que permite además llevar un control más efectivo de los saldos de los DTEs cuando se realizan ajustes como notas de crédito o débito.
type DTEBalanceControl struct {
	ID                            uint      `gorm:"primaryKey;autoIncrement:true;not null;index:idx_dte_balance_control"`
	BranchID                      uint      `gorm:"column:branch_id;type:uint;not null;index:idx_dte_branch"`
	OriginalDTEID                 string    `gorm:"column:original_dte_id;type:varchar(36);not null;index:idx_dte_original"`
	OriginalTaxedAmount           float64   `gorm:"column:original_taxed_amount;type:decimal(18,2);not null;"`
	OriginalExemptAmount          float64   `gorm:"column:original_exempt_amount;type:decimal(18,2);not null"`
	OriginalTotalNotSubjectAmount float64   `gorm:"column:original_not_subject_amount;type:decimal(18,2);not null"`
	RemainingTaxedAmount          float64   `gorm:"column:remaining_taxed_amount;type:decimal(18,2);not null"`
	RemainingExemptAmount         float64   `gorm:"column:remaining_exempt_amount;type:decimal(18,2);not null"`
	RemainingNotSubjectAmount     float64   `gorm:"column:remaining_not_subject_amount;type:decimal(18,2);not null"`
	CreatedAt                     time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt                     time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	OriginalDTE  *DTEDetails             `gorm:"foreignKey:OriginalDTEID;references:ID"`
	Branch       *BranchOffice           `gorm:"foreignKey:BranchID;references:ID"`
	Transactions []DTEBalanceTransaction `gorm:"foreignKey:BalanceControlID;references:ID"`
}

func (DTEBalanceControl) TableName() string {
	return "dte_balance_control"
}
