package db_models

import "time"

// SystemAdmin representa la tabla que almacenará la información de los administradores del sistema.
type SystemAdmin struct {
	ID        uint      `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(255);not null;uniqueIndex:idx_admin_email"`
	Password  string    `gorm:"column:password;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (SystemAdmin) TableName() string {
	return "system_admins"
}
