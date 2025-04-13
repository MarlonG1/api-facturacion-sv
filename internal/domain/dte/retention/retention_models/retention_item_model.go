package retention_models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
)

type RetentionItem struct {
	Number          item.ItemNumber         // Número del ítem de retención
	DTEType         document.DTEType        // Tipo de documento electrónico relacionado
	DocumentType    document.OperationType  // Tipo de generación del documento (Fisico o Electronico)
	DocumentNumber  document.DocumentNumber // Numero del documento relacionado
	EmissionDate    temporal.EmissionDate   // Fecha de emisión del documento relacionado
	RetentionAmount financial.Amount        // Monto de retención aplicado
	ReceptionCodeMH document.RetentionCode  // Código de retención del Ministerio de Hacienda
	RetentionIVA    financial.Amount        // Monto de retención del IVA
	Description     string                  // Descripción del ítem de retención
}
