package ccf_models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"

type CreditFiscalDocument struct {
	*models.DTEDocument               // Hereda la base de DTE
	CreditItems         []CreditItem  // Items específicos de CCF
	CreditSummary       CreditSummary // Resumen específico de CCF
}
