package interfaces

// DTEDocumentGetter es una interfaz que define los métodos getter que deben ser implementados por un documento DTE
type DTEDocumentGetter interface {
	GetIdentification() Identification      // GetIdentification retorna la identificación del documento
	GetAppendix() []Appendix                // GetAppendix retorna los anexos del documento
	GetExtension() Extension                // GetExtension retorna la extensión del documento
	GetIssuer() Issuer                      // GetIssuer retorna el emisor del documento
	GetReceiver() Receiver                  // GetReceiver retorna el receptor del documento
	GetItems() []Item                       // GetItems retorna los items del documento
	GetSummary() Summary                    // GetSummary retorna el resumen del documento
	GetRelatedDocuments() []RelatedDocument // GetRelatedDocuments retorna los documentos relacionados
	GetOtherDocuments() []OtherDocuments    // GetOtherDocuments retorna los otros documentos
	GetThirdPartySale() ThirdPartySale      // GetThirdPartySale retorna la venta de terceros
}

// DTEDocumentSetter es una interfaz que define los métodos setter que deben ser implementados por un documento DTE
type DTEDocumentSetter interface {
	SetIdentification(identification Identification) error        // SetIdentification establece la identificación del documento
	SetAppendix(appendix []Appendix) error                        // SetAppendix establece los anexos del documento
	SetExtension(extension Extension) error                       // SetExtension establece la extensión del documento
	SetIssuer(issuer Issuer) error                                // SetIssuer establece el emisor del documento
	SetReceiver(receiver Receiver) error                          // SetReceiver establece el receptor del documento
	SetItems(items []Item) error                                  // SetItems establece los items del documento
	SetSummary(summary Summary) error                             // SetSummary establece el resumen del documento
	SetRelatedDocuments(relatedDocuments []RelatedDocument) error // SetRelatedDocuments establece los documentos relacionados
	SetOtherDocuments(otherDocuments []OtherDocuments) error      // SetOtherDocuments establece los otros documentos
	SetThirdPartySale(thirdPartySale ThirdPartySale) error        // SetThirdPartySale establece la venta de terceros
}

// DTEDocument es una interfaz que combina los getters y setters de DTEDocument
type DTEDocument interface {
	DTEDocumentGetter
	DTEDocumentSetter
	Validate() error // Validate valida el documento
}
