package db_models

import "time"

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

func (DTEDocument) TableName() string {
	return "dte_documents"
}
