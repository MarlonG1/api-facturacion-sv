package db_models

// DTEDetails es una estructura que representa los detalles de un DTE almacenados en la base de datos.
// Se utiliza para almacenar los detalles de un DTE en la base de datos y recuperarlos para su procesamiento.
// Los detalles de un DTE se almacenan en la base de datos para su posterior procesamiento y envío a Hacienda.
//
// Para más información de la estructura de un DTE ver: https://factura.gob.sv/informacion-tecnica-y-funcional/
// en la sección de "Documentos de Sistema de Transmisión DTE", documento: "3. Manual Funcional del Sistema de Transmisión"
// página 53 del documento PDF.
//
// Nota: El DTE no se almacena firmado, solo se almacena en formato JSON, formato previo a la firma, la firma es el propio DTE
// firmado con la llave privada del emisor en formato JWT.
//
// El campo de DTEType indica el tipo de DTE, para más información sobre los tipos de DTE ver:
// https://factura.gob.sv/informacion-tecnica-y-funcional/
// en la sección de "Documentos de Sistema de Transmisión DTE", documento: "2. Catálogos- Sistema de Transmisión"
// página 5 del documento PDF y revisar /internal/domain/dte/common/constants/dte_type.go
type DTEDetails struct {
	ID             string  `gorm:"column:id;varchar(36);primaryKey;not null"`
	DTEType        string  `gorm:"column:dte_type;varchar(2);not null;index:idx_dte_type"`
	ControlNumber  string  `gorm:"column:control_number;varchar(30);not null;index"`
	ReceptionStamp *string `gorm:"column:reception_stamp;varchar(40)"`
	Transmission   string  `gorm:"column:transmission;varchar(15);not null"`
	Status         string  `gorm:"column:status;varchar(15);not null;index"`
	JSONData       string  `gorm:"column:json_data;type:json;not null"`
}

func (DTEDetails) TableName() string {
	return "dte_details"
}
