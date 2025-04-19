package interfaces

import "time"

// Identification es una interfaz que define los métodos que deben ser implementados por un objeto de tipo Identification
type Identification interface {
	IdentificationGetter
	IdentificationSetter
}

type IdentificationGetter interface {
	GetVersion() int               // GetVersion retorna la versión del DTE (Documento Tributario Electrónico)
	GetAmbient() string            // GetAmbient retorna el ambiente en el que se está emitiendo el DTE 00 -> Pruebas, 01 -> Producción
	GetDTEType() string            // GetDTEType retorna el tipo de DTE que se está emitiendo
	GetControlNumber() string      // GetControlNumber retorna el número de control del DTE
	GetGenerationCode() string     // GetGenerationCode retorna el código de generación del DTE
	GetModelType() int             // GetModelType retorna el tipo de modelo del DTE
	GetOperationType() int         // GetOperationType retorna el tipo de operación del DTE
	GetEmissionDate() time.Time    // GetEmissionDate retorna la fecha de emisión del DTE
	GetEmissionTime() time.Time    // GetEmissionTime retorna la hora de emisión del DTE
	GetCurrency() string           // GetCurrency retorna la moneda en la que se está emitiendo el DTE
	GetContingencyType() *int      // GetContingencyType retorna el tipo de contingencia en la que se está emitiendo el DTE
	GetContingencyReason() *string // GetContingencyReason retorna la razón de contingencia en la que se está emitiendo el DTE
}

type IdentificationSetter interface {
	SetControlNumber(controlNumber string) error          // SetControlNumber establece el número de control del DTE
	GenerateCode() error                                  // SetGenerationCode establece el código de generación del DTE
	SetVersion(version int) error                         // SetVersion establece la versión del DTE (Documento Tributario Electrónico)
	SetAmbient(ambient string) error                      // SetAmbient establece el ambiente en el que se está emitiendo el DTE 00 -> Pruebas, 01 -> Producción
	SetDTEType(dteType string) error                      // SetDTEType establece el tipo de DTE que se está emitiendo
	SetModelType(modelType int) error                     // SetModelType establece el tipo de modelo del DTE
	SetOperationType(operationType int) error             // SetOperationType establece el tipo de operación del DTE
	SetEmissionDate(emissionDate time.Time) error         // SetEmissionDate establece la fecha de emisión del DTE
	SetEmissionTime(emissionTime time.Time) error         // SetEmissionTime establece la hora de emisión del DTE
	SetCurrency(currency string) error                    // SetCurrency establece la moneda en la que se está emitiendo el DTE
	SetContingencyType(contingencyType *int) error        // SetContingencyType establece el tipo de contingencia en la que se está emitiendo el DTE
	SetContingencyReason(contingencyReason *string) error // SetContingencyReason establece la razón de contingencia en la que se está emitiendo el DTE
}
