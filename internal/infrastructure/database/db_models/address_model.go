package db_models

// Address Representa la dirección de la casa matriz, sucursal o agencias de emisores de DTEs en la base de datos.
// Los datos como Municipality y Department son códigos de dos letras porque se refieren a los códigos de los departamentos
// y municipios de El Salvador exigidos por Hacienda, para más información sobre los códigos de municipios y departamentos
//
// de El Salvador ver: https://factura.gob.sv/informacion-tecnica-y-funcional/
// en la sección de "Documentos de Sistema de Transmisión DTE", documento: "2. Catálogos- Sistema de Transmisión"
// página 6-8 del documento PDF y revisar internal/domain/dte/common/value_objects/location/department.go
// e internal/domain/dte/common/value_objects/location/municipality.go
type Address struct {
	ID           uint   `gorm:"column:id;type:uint;primaryKey;autoIncrement;not null"`
	BranchID     uint   `gorm:"column:branch_id;type:uint;not null;index:idx_address_branch"`
	Municipality string `gorm:"column:municipality;type:varchar(2);not null"`
	Department   string `gorm:"column:department;type:varchar(2);not null"`
	Complement   string `gorm:"column:complement;type:varchar(200);not null"`

	// Relaciones
	Branch *BranchOffice `gorm:"foreignKey:BranchID;references:ID"`
}

func (Address) TableName() string {
	return "addresses"
}
