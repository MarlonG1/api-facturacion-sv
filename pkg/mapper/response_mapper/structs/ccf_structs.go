package structs

type CCFDTEResponse struct {
	Identificacion       *DTEIdentification   `json:"identificacion"`
	Emisor               DTEIssuer            `json:"emisor"`
	Receptor             DTEReceiver          `json:"receptor"`
	CuerpoDocumento      []DTEItem            `json:"cuerpoDocumento"`
	Resumen              *DTESummary          `json:"resumen"`
	DocumentoRelacionado []DTERelatedDocument `json:"documentoRelacionado"`
	OtrosDocumentos      []DTEOtherDocument   `json:"otrosDocumentos"`
	VentaTercero         *DTEThirdPartySale   `json:"ventaTercero"`
	Extension            *DTEExtension        `json:"extension"`
	Apendice             []DTEApendice        `json:"apendice"`
}
