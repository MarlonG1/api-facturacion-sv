package db_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"gorm.io/gorm"
	"time"
)

// DTEBalanceTransaction representa una transacción de saldo de un DTE (Documento Tributario Electrónico).
// Su propósito es almacenar las transacciones de saldo de un DTE en la base de datos y recuperarlas para su procesamiento
// en esta tabla se registran las transacciones de Notas de Crédito y Débito que afectan el saldo de un DTE.
type DTEBalanceTransaction struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement:true;not null;index:idx_dte_balance_transaction"`
	BalanceControlID     uint      `gorm:"column:balance_control_id;type:int;not null;index:idx_dte_balance_control"`
	AdjustmentDocumentID string    `gorm:"column:adjustment_document_id;type:varchar(36);not null;index:idx_dte_adjustment_document"`
	TransactionType      string    `gorm:"column:transaction_type;type:varchar(2);not null;index:idx_dte_transaction_type"`
	TaxedAmount          float64   `gorm:"column:taxed_amount;type:decimal(18,2);not null"`
	ExemptAmount         float64   `gorm:"column:exempt_amount;type:decimal(18,2);not null"`
	NotSubjectAmount     float64   `gorm:"column:not_subject_amount;type:decimal(18,2);not null"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	BalanceControl     *DTEBalanceControl `gorm:"foreignKey:BalanceControlID;references:ID"`
	AdjustmentDocument *DTEDetails        `gorm:"foreignKey:AdjustmentDocumentID;references:ID"`
}

func (DTEBalanceTransaction) TableName() string {
	return "dte_balance_transactions"
}

func (d *DTEBalanceTransaction) AfterCreate(tx *gorm.DB) error {
	// Realizar una resta o suma en el balance de acuerdo al tipo de transacción
	if d.TransactionType == constants.NotaCreditoElectronica {
		d.BalanceControl.RemainingTaxedAmount -= d.TaxedAmount
		d.BalanceControl.RemainingExemptAmount -= d.ExemptAmount
		d.BalanceControl.RemainingNotSubjectAmount -= d.NotSubjectAmount
	} else {
		d.BalanceControl.RemainingTaxedAmount += d.TaxedAmount
		d.BalanceControl.RemainingExemptAmount += d.ExemptAmount
		d.BalanceControl.RemainingNotSubjectAmount += d.NotSubjectAmount
	}

	err := tx.Save(d.BalanceControl).Error
	if err != nil {
		return err
	}

	return nil
}
