package models

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
)

// Identification es una estructura que representa la identificación de un DTE, contiene Version, Ambient, DTEType, ControlNumber,
// GenerationCode, ModelType, OperationType, EmissionDate, EmissionTime, Currency, ContingencyType y ContingencyReason
type Identification struct {
	Version           document.Version              `json:"version"`
	Ambient           document.Ambient              `json:"ambient"`
	DTEType           document.DTEType              `json:"dteType"`
	ControlNumber     identification.ControlNumber  `json:"controlNumber"`
	GenerationCode    identification.GenerationCode `json:"generationCode"`
	ModelType         document.ModelType            `json:"modelType"`
	OperationType     document.OperationType        `json:"operationType"`
	EmissionDate      temporal.EmissionDate         `json:"emissionDate"`
	EmissionTime      temporal.EmissionTime         `json:"emissionTime"`
	Currency          financial.Currency            `json:"currency"`
	ContingencyType   *document.ContingencyType     `json:"contingencyType,omitempty"`
	ContingencyReason *document.ContingencyReason   `json:"contingencyReason,omitempty"`
}

func (i *Identification) GetVersion() int {
	return i.Version.GetValue()
}
func (i *Identification) GetAmbient() string {
	return i.Ambient.GetValue()
}
func (i *Identification) GetDTEType() string {
	return i.DTEType.GetValue()
}
func (i *Identification) GetControlNumber() string {
	return i.ControlNumber.GetValue()
}
func (i *Identification) GetGenerationCode() string {
	return i.GenerationCode.GetValue()
}
func (i *Identification) GetModelType() int {
	return i.ModelType.GetValue()
}
func (i *Identification) GetOperationType() int {
	return i.OperationType.GetValue()
}
func (i *Identification) GetEmissionDate() time.Time {
	return i.EmissionDate.GetValue()
}
func (i *Identification) GetEmissionTime() time.Time {
	return i.EmissionTime.GetValue()
}
func (i *Identification) GetCurrency() string {
	return i.Currency.GetValue()
}
func (i *Identification) GetContingencyType() *int {
	return utils.ToIntPointer(i.ContingencyType.GetValue())
}
func (i *Identification) GetContingencyReason() *string {
	return utils.ToStringPointer(i.ContingencyReason.GetValue())
}
func (i *Identification) SetControlNumber(controlNumber string) error {
	cn, err := identification.NewControlNumber(controlNumber)
	if err != nil {
		return err
	}
	i.ControlNumber = *cn
	return err
}
func (i *Identification) GenerateCode() error {
	gc, err := identification.NewGenerationCode()
	if err != nil {
		return err
	}
	i.GenerationCode = *gc
	return err
}
