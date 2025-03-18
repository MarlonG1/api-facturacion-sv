package interfaces

// ThirdPartySale es una interfaz que define los métodos que deben ser implementados por un objeto de tipo ThirdPartySale
type ThirdPartySale interface {
	GetNIT() string  // GetNIT retorna el NIT del tercero
	GetName() string // GetName retorna el nombre del tercero
}
