package interfaces

// ThirdPartySaleGetter es una interfaz que define los métodos getter que deben ser implementados por un objeto de tipo ThirdPartySale
type ThirdPartySaleGetter interface {
	GetNIT() string  // GetNIT retorna el NIT del tercero
	GetName() string // GetName retorna el nombre del tercero
}

// ThirdPartySaleSetter es una interfaz que define los métodos setter que deben ser implementados por un objeto de tipo ThirdPartySale
type ThirdPartySaleSetter interface {
	SetNIT(nit string) error   // SetNIT establece el NIT del tercero
	SetName(name string) error // SetName establece el nombre del tercero
}

// ThirdPartySale es una interfaz que combina los getters y setters de ThirdPartySale
type ThirdPartySale interface {
	ThirdPartySaleGetter
	ThirdPartySaleSetter
}
