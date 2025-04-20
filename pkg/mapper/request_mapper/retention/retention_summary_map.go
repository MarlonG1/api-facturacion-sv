package retention

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

func MapRetentionSummary(req *structs.RetentionSummary) (*retention_models.RetentionSummary, error) {
	if req == nil {
		return nil, dte_errors.NewValidationError("RequiredField", "RetentionSummary")
	}

	totalRetention, err := financial.NewAmountForTotal(req.TotalRetentionAmount)
	if err != nil {
		return nil, err
	}

	totalIVARetention, err := financial.NewAmountForTotal(req.TotalRetentionIVA)
	if err != nil {
		return nil, err
	}

	return &retention_models.RetentionSummary{
		TotalSubjectRetention: *totalRetention,
		TotalIVARetention:     *totalIVARetention,
	}, nil
}
