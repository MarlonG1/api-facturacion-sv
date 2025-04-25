package db_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"gorm.io/gorm"
	"time"
)

// DTEDocument representa la relación entre un documento tributario electrónico y una sucursal.
// Se utiliza para almacenar la relación entre un documento tributario electrónico y una sucursal.
// La relación entre un documento tributario electrónico y una sucursal se almacena en la base de datos
// para su posterior procesamiento y envío a Hacienda.
type DTEDocument struct {
	DocumentID string    `gorm:"column:document_id;type:varchar(36);primaryKey;not null;index:idx_dte_document"`
	BranchID   uint      `gorm:"column:branch_id;type:uint;not null;index:idx_dte_branch"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_dte_date"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Índice compuesto para consultas por período
	// `gorm:"index:idx_branch_period,priority:1,2"` - Para BranchID y CreatedAt

	// Relaciones
	Branch   *BranchOffice `gorm:"foreignKey:BranchID;references:ID"`
	Document *DTEDetails   `gorm:"foreignKey:DocumentID;references:ID"`
}

func (d *DTEDocument) AfterCreate(tx *gorm.DB) error {
	// Si el DTE es una nota de crédito o débito, se crea un registro en la tabla de control de saldo.
	if constants.ValidAdjustmentDTETypes[d.Document.DTEType] {

		extractor, err := utils.ExtractSummaryTotalAmountsFromStringJSON(d.Document.JSONData)
		if err != nil {
			return err
		}

		dteBalanceControl := &DTEBalanceControl{
			OriginalDTEID:                 d.DocumentID,
			BranchID:                      d.BranchID,
			OriginalTaxedAmount:           extractor.Summary.TotalTaxed,
			OriginalExemptAmount:          extractor.Summary.TotalExempt,
			OriginalTotalNotSubjectAmount: extractor.Summary.TotalNotSubject,
			RemainingTaxedAmount:          extractor.Summary.TotalTaxed,
			RemainingExemptAmount:         extractor.Summary.TotalExempt,
			RemainingNotSubjectAmount:     extractor.Summary.TotalNotSubject,
		}
		if err = tx.Create(dteBalanceControl).Error; err != nil {
			logs.Error("Error creating DTE balance control", map[string]interface{}{
				"error":       err.Error(),
				"document_id": d.DocumentID,
				"branch_id":   d.BranchID,
				"dteType":     d.Document.DTEType,
			})
			return err
		}
	}

	return nil
}

func (DTEDocument) TableName() string {
	return "dte_documents"
}
