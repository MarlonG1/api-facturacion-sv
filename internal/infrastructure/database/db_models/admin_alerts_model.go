package db_models

import "time"

// AdminAlert representa la estructura de la tabla admin_alerts en la base de datos.
// Esta tabla almacena los alertas que se envían a los administradores de la aplicación.
// Los alertas son generados por eventos de la aplicación y pueden ser de diferentes tipos y severidades.
// Las alertas generadas son enviadas a los administradores de la aplicación para que puedan tomar las acciones
// necesarias, dependiendo del tipo y severidad del alerta.
type AdminAlert struct {
	ID             uint      `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	EventID        uint      `gorm:"column:event_id;type:uint;not null;index:idx_admin_alerts_event"`
	AdminID        uint      `gorm:"column:admin_id;type:uint;not null;index:idx_admin_alerts_admin"`
	AlertType      string    `gorm:"column:alert_type;type:varchar(20);not null;index:idx_admin_alerts_alert_type"`
	Severity       string    `gorm:"column:severity;type:varchar(10);not null;index"`
	Message        string    `gorm:"column:message;type:text;not null"`
	DeliveryStatus string    `gorm:"column:delivery_status;type:varchar(20);not null;index"`
	DeliveryAt     time.Time `gorm:"column:delivery_at;type:datetime"`

	// Relaciones
	Admin *SystemAdmin `gorm:"foreignKey:AdminID;references:ID"`
	Event *DomainEvent `gorm:"foreignKey:EventID;references:ID"`
}

func (AdminAlert) TableName() string {
	return "admin_alerts"
}
