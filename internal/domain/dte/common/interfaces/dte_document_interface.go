package interfaces

// DTEDocument es una interfaz que define los métodos que deben ser implementados por un documento DTE
type DTEDocument interface {
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
	Validate() error                        // Validate valida el documento
}
