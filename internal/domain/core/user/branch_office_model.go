package user

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
)

type BranchOffice struct {
	ID                  uint     `json:"-"`
	UserID              uint     `json:"-"`
	EstablishmentCode   *string  `json:"establishment_code,omitempty"`
	Email               *string  `json:"email,omitempty"`
	APIKey              string   `json:"-"`
	APISecret           string   `json:"-"`
	Phone               *string  `json:"phone,omitempty"`
	EstablishmentType   string   `json:"establishment_type"`
	EstablishmentTypeMH *string  `json:"establishment_type_mh,omitempty"`
	POSCode             *string  `json:"pos_code,omitempty"`
	POSCodeMH           *string  `json:"pos_code_mh,omitempty"`
	IsActive            bool     `json:"is_active"`
	Address             *Address `json:"address,omitempty"`
}

func (b *BranchOffice) Validate() error {
	if b.EstablishmentCode != nil {
		if len(*b.EstablishmentCode) != 4 {
			return dte_errors.NewValidationError("InvalidLength", "establishment_code", "4", len(*b.EstablishmentCode))
		}
	}

	if b.Email != nil {
		if _, err := base.NewEmail(*b.Email); err != nil {
			return err
		}
	}

	if b.Phone != nil {
		if _, err := base.NewPhone(*b.Phone); err != nil {
			return err
		}
	}

	if b.EstablishmentTypeMH != nil {
		if _, err := document.NewEstablishmentType(*b.EstablishmentTypeMH); err != nil {
			return err
		}
	}

	if b.POSCode != nil {
		if len(*b.POSCode) != 4 {
			return dte_errors.NewValidationError("InvalidLength", "pos_code", "4", len(*b.POSCode))
		}
	}

	if b.POSCodeMH != nil {
		if len(*b.POSCodeMH) != 4 {
			return dte_errors.NewValidationError("InvalidLength", "pos_code_mh", "4", len(*b.POSCodeMH))
		}
	}

	if _, err := document.NewEstablishmentType(b.EstablishmentType); err != nil {
		return err
	}

	if b.Address.Validate() != nil {
		return b.Address.Validate()
	}

	return nil
}
