package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// Receiver es una estructura que representa el receptor de un DTE, contiene Name, DocumentType, DocumentNumber, Address,
// Email, Phone, NRC, ActivityDescription, ActivityCode y CommercialName
type Receiver struct {
	Name                *string                        `json:"name,omitempty"`
	DocumentType        *document.DTEType              `json:"documentType,omitempty"`
	DocumentNumber      *identification.DocumentNumber `json:"documentNumber,omitempty"`
	Address             interfaces.Address             `json:"address,omitempty"`
	Email               *base.Email                    `json:"email,omitempty"`
	Phone               *base.Phone                    `json:"phone,omitempty"`
	NRC                 *identification.NRC            `json:"nrc,omitempty"`
	NIT                 *identification.NIT            `json:"nit,omitempty"`
	ActivityDescription *string                        `json:"activityDescription,omitempty"`
	ActivityCode        *identification.ActivityCode   `json:"activityCode,omitempty"`
	CommercialName      *string                        `json:"commercialName,omitempty"`
}

func (r *Receiver) GetName() *string {
	return r.Name
}
func (r *Receiver) GetDocumentType() *string {
	if r.DocumentType == nil {
		return nil
	}

	return utils.ToStringPointer(r.DocumentType.GetValue())
}
func (r *Receiver) GetDocumentNumber() *string {
	if r.DocumentNumber == nil {
		return nil
	}

	return utils.ToStringPointer(r.DocumentNumber.GetValue())
}
func (r *Receiver) GetAddress() interfaces.Address {
	return r.Address
}
func (r *Receiver) GetEmail() *string {
	if r.Email == nil {
		return nil
	}

	return utils.ToStringPointer(r.Email.GetValue())
}
func (r *Receiver) GetPhone() *string {
	if r.Phone == nil {
		return nil
	}

	return utils.ToStringPointer(r.Phone.GetValue())
}
func (r *Receiver) GetNRC() *string {
	if r.NRC == nil {
		return nil
	}
	return utils.ToStringPointer(r.NRC.GetValue())
}
func (r *Receiver) GetActivityCode() *string {
	if r.ActivityCode == nil {
		return nil
	}
	return utils.ToStringPointer(r.ActivityCode.GetValue())
}
func (r *Receiver) GetNIT() *string {
	if r.NIT == nil {
		return nil
	}
	return utils.ToStringPointer(r.NIT.GetValue())
}
func (r *Receiver) GetActivityDescription() *string {
	return r.ActivityDescription
}
func (r *Receiver) GetCommercialName() *string {
	return r.CommercialName
}
