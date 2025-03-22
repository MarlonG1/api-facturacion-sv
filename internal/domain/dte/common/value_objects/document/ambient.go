package document

import (
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type Ambient struct {
	Value string `json:"value"`
}

// NewAmbient Crea un nuevo objeto de valor Ambient con el valor del ambiente obtenido del entorno
func NewAmbient() (*Ambient, error) {
	ambient := &Ambient{
		Value: env.Server.AmbientCode,
	}

	if ambient.IsValid() {
		return ambient, nil
	} else {
		return &Ambient{}, dte_errors.NewValidationError("InvalidAmbientCode", ambient.Value)
	}
}

func NewValidatedAmbient(value string) *Ambient {
	return &Ambient{Value: value}
}

func NewAmbientCustom(value string) (*Ambient, error) {
	ambient := &Ambient{Value: value}
	if ambient.IsValid() {
		return ambient, nil
	} else {
		return &Ambient{}, dte_errors.NewValidationError("InvalidAmbientCode", ambient.Value)
	}
}

/*
IsValid Válida que el valor del campo Ambient sea uno de los valores permitidos y en caso contrario lanza un error de validación
Si un error de validación es lanzado, este se propaga a través de la infraestructura de errores
*/
func (a *Ambient) IsValid() bool {
	for _, v := range constants.AllowedAmbientValues {
		if a.Value == v {
			return true
		}
	}
	return false
}

// Equals Compara el valor del campo Ambient con el valor de otro objeto de valor Ambient
func (a *Ambient) Equals(ambient interfaces.ValueObject[string]) bool {
	return a.GetValue() == ambient.GetValue()
}

func (a *Ambient) GetValue() string {
	return a.Value
}

func (a *Ambient) ToString() string { return a.Value }
