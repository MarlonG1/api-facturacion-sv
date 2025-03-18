package db_models

// DomainEvent representa un evento de dominio que se guarda en la base de datos
// Su función es guardar los eventos de dominio que se generan en la aplicación por alguna situación específica que requiera
// la atención de los usuarios.
// Por ejemplo, cuando una sucursal entra en estado de contingencia, se genera un evento de dominio para notificar al usuario
type DomainEvent struct {
	ID         uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	UserID     uint   `gorm:"column:user_id;type:uint;not null;index:idx_event_user"`
	BranchID   uint   `gorm:"column:branch_id;type:uint;not null;index:idx_event_branch"`
	EventType  string `gorm:"column:event_type;type:varchar(50);not null;index"`
	Payload    string `gorm:"column:payload;type:json;not null"`
	OccurredAt string `gorm:"column:occurred_at;type:timestamp;not null;index"`

	// Índice compuesto
	// `gorm:"index:idx_user_type,priority:1,2"` - Para buscar eventos de un usuario por tipo

	// Relaciones
	User   *User         `gorm:"foreignKey:UserID;references:ID"`
	Branch *BranchOffice `gorm:"foreignKey:BranchID;references:ID"`
}

func (DomainEvent) TableName() string {
	return "domain_events"
}
