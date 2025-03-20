package db_models

import "time"

// ContingencyDocument representa un documento de contingencia en la base de datos
// Se utiliza para almacenar documentos que no pudieron ser enviados a Hacienda y se deben enviar de manera diferida
// La aplicación se encarga de reintentar el envío de estos documentos de manera automática.
// El campo Type indica el tipo de documento de contingencia, para más información sobre los tipos contingencia
//
// ver: https://factura.gob.sv/informacion-tecnica-y-funcional/
// en la sección de "Documentos de Sistema de Transmisión DTE", documento: "2. Catálogos- Sistema de Transmisión"
// página 5 del documento PDF y revisar /internal/domain/dte/common/constants/contingency_document_types.go
type ContingencyDocument struct {
	ID              string    `gorm:"column:id;type:varchar(36);primaryKey;not null"`
	DocumentID      string    `gorm:"column:document_id;type:varchar(36);not null;index"`
	BranchID        uint      `gorm:"column:branch_id;type:uint;not null;index:idx_contingency_branch"`
	ContingencyType int8      `gorm:"column:type;contingency_type:tinyint;not null;index"`
	Reason          string    `gorm:"column:reason;type:varchar(150);not null;index"`
	BatchID         *string   `gorm:"column:batch_id;type:varchar(36);index"`
	MHBatchID       *string   `gorm:"column:mh_batch_id;type:varchar(36)"`
	Observations    *string   `gorm:"column:observations;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_contingency_date"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Índice compuesto
	// `gorm:"index:idx_branch_date,priority:1,2"` - Este índice sería para BranchID y CreatedAt

	// Relaciones
	Document *DTEDetails   `gorm:"foreignKey:DocumentID;references:ID"`
	Branch   *BranchOffice `gorm:"foreignKey:BranchID;references:ID"`
}

func (ContingencyDocument) TableName() string {
	return "contingency_documents"
}
