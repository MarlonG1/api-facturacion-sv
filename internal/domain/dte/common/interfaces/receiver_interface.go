package interfaces

// Receiver es una interfaz que define los métodos que debe implementar un receptor
type Receiver interface {
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
