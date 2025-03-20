package db_models

import "time"

// User representa la tabla cuya información será crucial para el proceso de autenticación y autorización de los usuarios
// que desean utilizar la API de Factura Electrónica. Esta tabla almacenará la información de los usuarios que desean
// utilizar la API de Factura Electrónica.
//
// El campo PasswordPri la contraseña privada del certificado digital del usuario brindada por el Ministerio de Hacienda,
// esta contraseña es necesaria para firmar los documentos electrónicos.
//
// El campo YearInDTE será un indicativo para determinar si el usuario desea que el año de emisión del documento electrónico
// se muestre en el número de control. Por ejemplo:
//  1. Si es true, el número de control será: DTE-01-00000000-202500000000001
//  2. Si es false, el número de control será: DTE-01-00000000-000000000000001
type User struct {
	ID             uint      `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	NIT            string    `gorm:"column:nit;type:varchar(17);not null;uniqueIndex"`
	NRC            string    `gorm:"column:nrc;type:varchar(10);not null;uniqueIndex"`
	Status         bool      `gorm:"column:status;type:tinyint;not null;index:idx_user_status"`
	AuthType       string    `gorm:"column:auth_type;type:varchar(15);not null"`
	PasswordPri    string    `gorm:"column:password_pri;type:varchar(255);not null"`
	CommercialName string    `gorm:"column:commercial_name;type:varchar(255);not null"`
	Business       string    `gorm:"column:business_name;type:varchar(255);not null"`
	Email          string    `gorm:"column:email;type:varchar(255);not null;uniqueIndex:idx_user_email"`
	YearInDTE      bool      `gorm:"column:year_in_dte;type:tinyint;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "users"
}
