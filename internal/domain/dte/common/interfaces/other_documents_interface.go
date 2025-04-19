package interfaces

// OtherDocumentsGetter es una interfaz que define los métodos getter que debe implementar un documento
type OtherDocumentsGetter interface {
	GetAssociatedDocument() int // GetAssociatedDocument obtiene el documento asociado
	GetDescription() string     // GetDescription obtiene la descripción del documento
	GetDetail() string          // GetDetail obtiene el detalle del documento
	GetDoctor() DoctorInfo      // GetDoctor obtiene la información del doctor
}

// OtherDocumentsSetter es una interfaz que define los métodos setter que debe implementar un documento
type OtherDocumentsSetter interface {
	SetAssociatedDocument(associatedDocument int) error // SetAssociatedDocument establece el documento asociado
	SetDescription(description string) error            // SetDescription establece la descripción del documento
	SetDetail(detail string) error                      // SetDetail establece el detalle del documento
	SetDoctor(doctor DoctorInfo) error                  // SetDoctor establece la información del doctor
}

// OtherDocuments es una interfaz que combina los getters y setters de OtherDocuments
type OtherDocuments interface {
	OtherDocumentsGetter
	OtherDocumentsSetter
}

// DoctorInfoGetter es una interfaz que define los métodos getter que debe implementar un doctor
type DoctorInfoGetter interface {
	GetName() string           // GetName obtiene el nombre del doctor
	GetServiceType() int       // GetServiceType obtiene el tipo de servicio del doctor
	GetNIT() string            // GetNIT obtiene el NIT del doctor
	GetIdentification() string // GetIdentification obtiene la identificación del doctor
}

// DoctorInfoSetter es una interfaz que define los métodos setter que debe implementar un doctor
type DoctorInfoSetter interface {
	SetName(name string) error                     // SetName establece el nombre del doctor
	SetServiceType(serviceType int) error          // SetServiceType establece el tipo de servicio del doctor
	SetNIT(nit string) error                       // SetNIT establece el NIT del doctor
	SetIdentification(identification string) error // SetIdentification establece la identificación del doctor
}

// DoctorInfo es una interfaz que combina los getters y setters de DoctorInfo
type DoctorInfo interface {
	DoctorInfoGetter
	DoctorInfoSetter
}
