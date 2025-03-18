package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
)

// Issuer es una estructura que representa el emisor de un DTE
type Issuer struct {
	NIT                 identification.NIT          `json:"nit"`
	NRC                 identification.NRC          `json:"nrc"`
	Name                string                      `json:"name"`
	ActivityCode        identification.ActivityCode `json:"activityCode"`
	ActivityDescription string                      `json:"activityDescription"`
	EstablishmentType   document.EstablishmentType  `json:"establishmentType"`
	Address             interfaces.Address          `json:"address"`
	Phone               base.Phone                  `json:"phone"`
	Email               base.Email                  `json:"email"`
	CommercialName      string                      `json:"commercialName"`
	EstablishmentCode   *string                     `json:"establishmentCode,omitempty"`
	EstablishmentMHCode *string                     `json:"establishmentMHCode,omitempty"`
	POSCode             *string                     `json:"POSCode,omitempty"`
	POSMHCode           *string                     `json:"POSMHCode,omitempty"`
}

func (i *Issuer) GetName() string {
	return i.Name
}
func (i *Issuer) GetActivityDescription() string {
	return i.ActivityDescription
}
func (i *Issuer) GetCommercialName() string {
	return i.CommercialName
}
func (i *Issuer) GetNIT() string {
	return i.NIT.GetValue()
}
func (i *Issuer) GetNRC() string {
	return i.NRC.GetValue()
}
func (i *Issuer) GetActivityCode() string {
	return i.ActivityCode.GetValue()
}
func (i *Issuer) GetEstablishmentType() string {
	return i.EstablishmentType.GetValue()
}
func (i *Issuer) GetAddress() interfaces.Address {
	return i.Address
}
func (i *Issuer) GetPhone() string {
	return i.Phone.GetValue()
}
func (i *Issuer) GetEmail() string {
	return i.Email.GetValue()
}
func (i *Issuer) GetEstablishmentCode() *string {
	return i.EstablishmentCode
}
func (i *Issuer) GetEstablishmentMHCode() *string {
	return i.EstablishmentMHCode
}
func (i *Issuer) GetPOSCode() *string {
	return i.POSCode
}
func (i *Issuer) GetPOSMHCode() *string {
	return i.POSMHCode
}
