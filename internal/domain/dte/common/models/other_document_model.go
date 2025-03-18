package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type OtherDocument struct {
	AssociatedCode document.AssociatedDocumentCode `json:"associatedCode"`
	Description    *string                         `json:"description,omitempty"`
	Detail         *string                         `json:"detail,omitempty"`
	Doctor         interfaces.DoctorInfo           `json:"doctor,omitempty"`
}

type DoctorInfo struct {
	Name           string               `json:"name"`
	ServiceType    document.ServiceType `json:"serviceType"`
	NIT            *identification.NIT  `json:"NIT,omitempty"`
	Identification *string              `json:"identification,omitempty"`
}

func (o *OtherDocument) GetAssociatedDocument() int {
	return o.AssociatedCode.GetValue()
}

func (o *OtherDocument) GetDescription() string {
	return utils.PointerToString(o.Description)
}

func (o *OtherDocument) GetDetail() string {
	return utils.PointerToString(o.Detail)
}

func (o *OtherDocument) GetDoctor() interfaces.DoctorInfo {
	return o.Doctor
}

func (d *DoctorInfo) GetName() string {
	if d == nil {
		return ""
	}
	return d.Name
}

func (d *DoctorInfo) GetServiceType() int {
	if d == nil {
		return 0
	}
	return d.ServiceType.GetValue()
}

func (d *DoctorInfo) GetNIT() string {
	if d == nil || d.NIT == nil {
		return ""
	}
	return d.NIT.GetValue()
}

func (d *DoctorInfo) GetIdentification() string {
	if d == nil {
		return ""
	}
	return utils.PointerToString(d.Identification)
}
