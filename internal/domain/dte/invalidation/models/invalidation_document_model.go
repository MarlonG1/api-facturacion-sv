package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
)

type InvalidationDocument struct {
	Identification *models.Identification `json:"identification,omitempty"`
	Issuer         *models.Issuer         `json:"issuer,omitempty"`
	Document       *InvalidatedDocument   `json:"document,omitempty"`
	Reason         *InvalidationReason    `json:"reason,omitempty"`
}
