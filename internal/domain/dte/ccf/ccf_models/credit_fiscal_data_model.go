package ccf_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type CCFData struct {
	*models.InputDataCommon
	Items         []CreditItem
	CreditSummary *CreditSummary
}
