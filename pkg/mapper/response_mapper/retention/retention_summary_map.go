package retention

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func MapRetentionResponseSummary(summary *retention_models.RetentionSummary) *structs.RetentionSummary {
	if summary == nil {
		return nil
	}

	return &structs.RetentionSummary{
		TotalIvaRetenido:       summary.TotalIVARetention.GetValue(),
		TotalSujRetencion:      summary.TotalSubjectRetention.GetValue(),
		TotalIvaRetenidoLetras: summary.TotalIVARetentionLetters,
	}
}
