package interfaces

// IssuerGetter es una interfaz que define los métodos getter que debe implementar un emisor
type IssuerGetter interface {
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

// IssuerSetter es una interfaz que define los métodos setter que debe implementar un emisor
type IssuerSetter interface {
	SetName(name string) error                                // SetName establece el nombre del emisor
	SetActivityDescription(description string) error          // SetActivityDescription establece la descripción de la actividad económica del emisor
	SetCommercialName(commercialName string) error            // SetCommercialName establece el nombre comercial del emisor
	SetNIT(nit string) error                                  // SetNIT establece el NIT del emisor
	SetNRC(nrc string) error                                  // SetNRC establece el NRC del emisor
	SetActivityCode(activityCode string) error                // SetActivityCode establece el código de actividad económica del emisor
	SetEstablishmentType(establishmentType string) error      // SetEstablishmentType establece el tipo de establecimiento del emisor
	SetAddress(address Address) error                         // SetAddress establece la dirección del emisor
	SetPhone(phone string) error                              // SetPhone establece el teléfono del emisor
	SetEmail(email string) error                              // SetEmail establece el correo electrónico del emisor
	SetEstablishmentCode(establishmentCode *string) error     // SetEstablishmentCode establece el código de establecimiento del emisor
	SetEstablishmentMHCode(establishmentMHCode *string) error // SetEstablishmentMHCode establece el código de establecimiento de matriz o sucursal del emisor
	SetPOSCode(posCode *string) error                         // SetPOSCode establece el código de punto de venta del emisor
	SetPOSMHCode(posMHCode *string) error                     // SetPOSMHCode establece el código de punto de venta de matriz o sucursal del emisor
}

// Issuer es una interfaz que combina los getters y setters de Issuer
type Issuer interface {
	IssuerGetter
	IssuerSetter
}
