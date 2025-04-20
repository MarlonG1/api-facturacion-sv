package contingency

import "context"

// ContingencyManager interfaz para manejo de documentos en contingencia
type ContingencyManager interface {
	// StoreDocumentInContingency almacena un documento en contingencia
	StoreDocumentInContingency(ctx context.Context, document interface{}, dteType string, contingencyType int8, reason string) error
	// RetransmitPendingDocuments retransmite los documentos pendientes
	RetransmitPendingDocuments(ctx context.Context) error
}
