package interfaces

// AddressGetter es una interfaz que define los métodos getter que debe implementar una dirección
type AddressGetter interface {
	GetDepartment() string   // GetDepartment obtiene el departamento de la dirección
	GetMunicipality() string // GetMunicipality obtiene el municipio de la dirección
	GetComplement() string   // GetComplement obtiene el complemento de la dirección
}

// AddressSetter es una interfaz que define los métodos setter que debe implementar una dirección
type AddressSetter interface {
	SetDepartment(department string) error     // SetDepartment establece el departamento de la dirección
	SetMunicipality(municipality string) error // SetMunicipality establece el municipio de la dirección
	SetComplement(complement string) error     // SetComplement establece el complemento de la dirección
}

// Address es una interfaz que combina los getters y setters de Address
type Address interface {
	AddressGetter
	AddressSetter
}
