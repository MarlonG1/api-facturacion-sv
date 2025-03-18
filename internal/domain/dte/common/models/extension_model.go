package models

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
)

// Extension es una estructura que representa una extensión de un DTE, contiene DeliveryName, DeliveryDocument, ReceiverName,
// ReceiverDocument y Observation
type Extension struct {
	DeliveryName     document.DeliveryName     `json:"deliveryName"`
	DeliveryDocument document.DeliveryDocument `json:"deliveryDocument"`
	ReceiverName     document.DeliveryName     `json:"receiverName"`
	ReceiverDocument document.DeliveryDocument `json:"receiverDocument"`
	Observation      *document.Observation     `json:"observation,omitempty"`
	VehiculePlate    *string                   `json:"vehiculePlate,omitempty"`
}

func (e *Extension) GetDeliveryName() string {
	return e.DeliveryName.GetValue()
}
func (e *Extension) GetDeliveryDocument() string {
	return e.DeliveryDocument.GetValue()
}
func (e *Extension) GetReceiverName() string {
	return e.ReceiverName.GetValue()
}
func (e *Extension) GetReceiverDocument() string {
	return e.ReceiverDocument.GetValue()
}
func (e *Extension) GetVehiculePlate() *string {
	return e.VehiculePlate
}
func (e *Extension) GetObservation() *string {
	if e.Observation != nil {
		value := e.Observation.GetValue()
		return &value
	}
	return nil
}
