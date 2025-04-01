package constants

const (
	DocumentoEmisor = iota + 1
	DocumentoReceptor
	DocumentoMedico
	DocumentoTransporte
)

var (
	// AllowedAssociatedDocumentCodes contiene los tipos de documentos asociados permitidos, usado para validaciones
	AllowedAssociatedDocumentCodes = []int{
		DocumentoEmisor,
		DocumentoReceptor,
		DocumentoMedico,
		DocumentoTransporte,
	}
)
