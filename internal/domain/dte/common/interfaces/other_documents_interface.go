package interfaces

// OtherDocuments es una interfaz que define los métodos que debe implementar un documento
type OtherDocuments interface {
	GetAssociatedDocument() int // GetAssociatedDocument obtiene el documento asociado
	GetDescription() string     // GetDescription obtiene la descripción del documento
	GetDetail() string          // GetDetail obtiene el detalle del documento
	GetDoctor() DoctorInfo      // GetDoctor obtiene la información del doctor
}

// DoctorInfo es una interfaz que define los métodos que debe implementar un doctor
type DoctorInfo interface {
	GetName() string           // GetName obtiene el nombre del doctor
	GetServiceType() int       // GetServiceType obtiene el tipo de servicio del doctor
	GetNIT() string            // GetNIT obtiene el NIT del doctor
	GetIdentification() string // GetIdentification obtiene la identificación del doctor
}
