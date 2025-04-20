package interfaces

// ExtensionGetter es una interfaz que define los métodos getter que deben ser implementados por una extensión
type ExtensionGetter interface {
	GetDeliveryName() string     // GetDeliveryName obtiene el nombre del destinatario
	GetDeliveryDocument() string // GetDeliveryDocument obtiene el documento del destinatario
	GetReceiverName() string     // GetReceiverName obtiene el nombre del receptor
	GetReceiverDocument() string // GetReceiverDocument obtiene el documento del receptor
	GetObservation() *string     // GetObservation obtiene la observación
	GetVehiculePlate() *string   // GetVehiculePlate obtiene la placa del vehículo
}

// ExtensionSetter es una interfaz que define los métodos setter que deben ser implementados por una extensión
type ExtensionSetter interface {
	SetDeliveryName(deliveryName string) error         // SetDeliveryName establece el nombre del destinatario
	SetDeliveryDocument(deliveryDocument string) error // SetDeliveryDocument establece el documento del destinatario
	SetReceiverName(receiverName string) error         // SetReceiverName establece el nombre del receptor
	SetReceiverDocument(receiverDocument string) error // SetReceiverDocument establece el documento del receptor
	SetObservation(observation *string) error          // SetObservation establece la observación
	SetVehiculePlate(vehiculePlate *string) error      // SetVehiculePlate establece la placa del vehículo
}

// Extension es una interfaz que combina los getters y setters de Extension
type Extension interface {
	ExtensionGetter
	ExtensionSetter
}
