package retention_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"

type RetentionSummary struct {
	TotalSubjectRetention    financial.Amount
	TotalIVARetention        financial.Amount
	TotalIVARetentionLetters string
}
