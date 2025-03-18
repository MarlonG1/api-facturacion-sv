package db_models

import "time"

// AdminAlertRecipients representa la tabla admin_alert_recipients en la base de datos.
// Cada fila representa una relación entre un alerta y un administrador que recibirá la alerta,
// su funcionalidad es la de retransmitir alertas a los administradores de la aplicación
// las alertas son generadas por eventos en la aplicación que requieren la atención de un administrador
type AdminAlertRecipients struct {
	ID        uint      `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	AlertID   uint      `gorm:"column:alert_id;type:uint;not null;index:idx_alert_recipients_alert"`
	AdminID   uint      `gorm:"column:admin_id;type:uint;not null;index:idx_alert_recipients_admin"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relaciones
	Admin *SystemAdmin `gorm:"foreignKey:AdminID;references:ID"`
	Alert *AdminAlert  `gorm:"foreignKey:AlertID;references:ID"`
}

func (AdminAlertRecipients) TableName() string {
	return "admin_alert_recipients"
}
