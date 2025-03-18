package interfaces

// Address es una interfaz que define los métodos que debe implementar una dirección
type Address interface {
	GetDepartment() string   // GetDepartment obtiene el departamento de la dirección
	GetMunicipality() string // GetMunicipality obtiene el municipio de la dirección
	GetComplement() string   // GetComplement obtiene el complemento de la dirección
}
