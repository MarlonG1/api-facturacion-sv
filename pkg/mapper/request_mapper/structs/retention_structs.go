package structs

/*
	Para tipo documento 1 (Fisico) solicito:
		- Tipo de documento
		- Numero de documento (Numero correlativo tradicional)
		- Descripcion
		- Monto gravado
		- Fecha de Emision
		- Tipo de DTE
		- Codigo de retencion de MH
		- Toda la seccion de "receptor" es obligatoria

	Para tipo documento 2 (Electronico) solicito:
		- Descripcion
		- Tipo de documento
		- Numero de documento (Codigo de generacion UUID)
		- Codigo de retencion de MH (Sigo considerarlo si interpretarlo automaticamente...)
		- Toda la sección del "receptor" queda excluida


	NOTAS:
		- Los campos Extension y Appendixes son opcionales en ambos casos
		- Los campos "taxed_amount", "emission_date" los calculare automaticamente si el tipo documento es electronico, los tomo de la DB
 		  de los DTE previamente emitidos
*/

type RetentionItem struct {
	DocumentType   int      `json:"type"`
	DocumentNumber string   `json:"document_number"`
	Description    string   `json:"description"`
	RetentionCode  string   `json:"retention_code"`
	IvaAmount      *float64 `json:"iva_amount,omitempty"`
	TaxedAmount    *float64 `json:"taxed_amount,omitempty"`
	EmissionDate   *string  `json:"emission_date,omitempty"`
	DTEType        *string  `json:"dte_type,omitempty"`
}

type RetentionSummary struct {
	TotalRetentionAmount float64 `json:"total_retention_amount"`
	TotalRetentionIVA    float64 `json:"total_retention_iva"`
}

type CreateRetentionRequest struct {
	Items      []RetentionItem   `json:"items"`
	Summary    *RetentionSummary `json:"summary,omitempty"`
	Receiver   *ReceiverRequest  `json:"receiver,omitempty"`
	Extension  *ExtensionRequest `json:"extension,omitempty"`
	Appendixes []AppendixRequest `json:"appendixes,omitempty"`
}
