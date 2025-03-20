package db_models

// NotifiableUser representa la tabla cuya información será crucial para el proceso de notificaciones que generaran mediante
// los eventos de dominio. Esta tabla almacenará la información de los usuarios que desean recibir notificaciones.
// EntityID y EntityType representa que dichos atributos son polimórficos, es decir, pueden ser un Usuario Cliente o un
// Usuario Administrador.
type NotifiableUser struct {
	ID          uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	UserID      uint   `gorm:"column:user_id;type:uint;not null;index:idx_entity,priority:1"`
	EntityType  string `gorm:"column:entity_type;type:varchar(10);not null;index:idx_entity,priority:2"`
	Email       string `gorm:"column:email;type:varchar(255);not null;index:idx_notifiable_email"`
	EnabledPush bool   `gorm:"column:enabled_push;type:tinyint;not null;index:idx_push_enabled"`
}

func (NotifiableUser) TableName() string {
	return "notification_users"
}
