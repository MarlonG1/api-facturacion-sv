package interfaces

// ReceiverGetter es una interfaz que define los métodos getter que debe implementar un receptor
type ReceiverGetter interface {
	GetName() *string                // GetName obtiene el nombre del receptor
	GetDocumentType() *string        // GetDocumentType obtiene el tipo de documento del receptor
	GetDocumentNumber() *string      // GetDocumentNumber obtiene el número de documento del receptor
	GetAddress() Address             // GetAddress obtiene la dirección del receptor
	GetEmail() *string               // GetEmail obtiene el correo electrónico del receptor
	GetPhone() *string               // GetPhone obtiene el teléfono del receptor
	GetNRC() *string                 // GetNRC obtiene el NRC del receptor
	GetNIT() *string                 // GetNIT obtiene el NIT del receptor
	GetActivityCode() *string        // GetActivityCode obtiene el código de actividad económica del receptor
	GetActivityDescription() *string // GetActivityDescription obtiene la descripción de la actividad económica del receptor
	GetCommercialName() *string      // GetCommercialName obtiene el nombre comercial del receptor
}

// ReceiverSetter es una interfaz que define los métodos setter que debe implementar un receptor
type ReceiverSetter interface {
	SetName(name *string) error                               // SetName establece el nombre del receptor
	SetDocumentType(documentType *string) error               // SetDocumentType establece el tipo de documento del receptor
	SetDocumentNumber(documentNumber *string) error           // SetDocumentNumber establece el número de documento del receptor
	SetAddress(address Address) error                         // SetAddress establece la dirección del receptor
	SetEmail(email *string) error                             // SetEmail establece el correo electrónico del receptor
	SetPhone(phone *string) error                             // SetPhone establece el teléfono del receptor
	SetNRC(nrc *string) error                                 // SetNRC establece el NRC del receptor
	SetNIT(nit *string) error                                 // SetNIT establece el NIT del receptor
	SetActivityCode(activityCode *string) error               // SetActivityCode establece el código de actividad económica del receptor
	SetActivityDescription(activityDescription *string) error // SetActivityDescription establece la descripción de la actividad económica del receptor
	SetCommercialName(commercialName *string) error           // SetCommercialName establece el nombre comercial del receptor
}

// Receiver es una interfaz que combina los getters y setters de Receiver
type Receiver interface {
	ReceiverGetter
	ReceiverSetter
}
