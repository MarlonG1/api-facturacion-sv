package interfaces

// Issuer es una interfaz que define los métodos que debe implementar un emisor de un documento tributario electrónico
type Issuer interface {
	GetName() string                 // GetName retorna el nombre del emisor
	GetActivityDescription() string  // GetActivityDescription retorna la descripción de la actividad económica del emisor
	GetCommercialName() string       // GetCommercialName retorna el nombre comercial del emisor
	GetNIT() string                  // GetNIT retorna el NIT del emisor
	GetNRC() string                  // GetNRC retorna el NRC del emisor
	GetActivityCode() string         // GetActivityCode retorna el código de actividad económica del emisor
	GetEstablishmentType() string    // GetEstablishmentType retorna el tipo de establecimiento del emisor
	GetAddress() Address             // GetAddress retorna la dirección del emisor
	GetPhone() string                // GetPhone retorna el teléfono del emisor
	GetEmail() string                // GetEmail retorna el correo electrónico del emisor
	GetEstablishmentCode() *string   // GetEstablishmentCode retorna el código de establecimiento del emisor
	GetEstablishmentMHCode() *string // GetEstablishmentMHCode retorna el código de establecimiento de matriz o sucursal del emisor
	GetPOSCode() *string             // GetPOSCode retorna el código de punto de venta del emisor
	GetPOSMHCode() *string           // GetPOSMHCode retorna el código de punto de venta de matriz o sucursal del emisor
}
