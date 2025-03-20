package db_models

// BranchOffice representa la estructura de la tabla branch_offices en la base de datos
// El campo EstablishmentType es un campo de 2 caracteres que representa el tipo de establecimiento exigido por la Hacienda
// el campo hace referencia a si lugar es una sucursal, casa matriz, etc.
//
// Para mas informacion sobre los tipos de establecimientos ver:
// https://factura.gob.sv/informacion-tecnica-y-funcional/
// en la sección de "Documentos de Sistema de Transmisión DTE", documento: "2. Catálogos- Sistema de Transmisión"
// página 6 del documento PDF y revisar /internal/domain/dte/common/constants/establishment_type.go
//
// Los campos con terminación MH hacen referencia a los códigos brindados por Hacienda, si no posee dichos códigos de Hacienda
// dejar los campos en blanco para evitar problemas legales
type BranchOffice struct {
	ID                  uint    `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	UserID              uint    `gorm:"column:user_id;type:uint;not null;index:idx_branch_offices_user"`
	EstablishmentCode   *string `gorm:"column:establishment_code;type:varchar(4)"`
	Email               *string `gorm:"column:email;type:varchar(255)"`
	APIKey              string  `gorm:"column:api_key;type:varchar(255);not null;uniqueIndex"`
	APISecret           string  `gorm:"column:api_secret;type:varchar(255);not null"`
	Phone               *string `gorm:"column:phone;type:varchar(8)"`
	EstablishmentType   string  `gorm:"column:establishment_type;type:varchar(2);not null;index:idx_branch_est_type"`
	EstablishmentTypeMH *string `gorm:"column:establishment_type_mh;type:varchar(4)"`
	POSCode             *string `gorm:"column:pos_code;type:varchar(4)"`
	POSCodeMH           *string `gorm:"column:pos_code_mh;type:varchar(4)"`
	IsActive            bool    `gorm:"column:is_active;type:tinyint(1);not null;index:idx_branch_offices_active"`

	// Relaciones
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (BranchOffice) TableName() string {
	return "branch_offices"
}
