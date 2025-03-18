package interfaces

// Extension es una interfaz que define los métodos que deben ser implementados por una extensión
type Extension interface {
	GetDeliveryName() string     // GetDeliveryName obtiene el nombre del destinatario
	GetDeliveryDocument() string // GetDeliveryDocument obtiene el documento del destinatario
	GetReceiverName() string     // GetReceiverName obtiene el nombre del receptor
	GetReceiverDocument() string // GetReceiverDocument obtiene el documento del receptor
	GetObservation() *string     // GetObservation obtiene la observación
	GetVehiculePlate() *string   // GetVehiculePlate obtiene la placa del vehículo
}
