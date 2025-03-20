package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
)

type InvalidationReason struct {
	Type               document.InvalidationType     `json:"type"`
	ResponsibleName    string                        `json:"responsibleName"`
	ResponsibleDocType document.DTEType              `json:"responsibleDocType"`
	ResponsibleDocNum  identification.DocumentNumber `json:"responsibleDocNum"`
	RequesterName      string                        `json:"requesterName"`
	RequesterDocType   document.DTEType              `json:"requesterDocType"`
	RequesterDocNum    identification.DocumentNumber `json:"requesterDocNum"`
	Reason             *document.InvalidationReason  `json:"reason"` // null si no es tipo 3
}
