package credit_note

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func MapCreditNoteResponseExtension(extension interfaces.Extension) *structs.CreditNoteDTEExtension {
	if extension == nil {
		return nil
	}

	return &structs.CreditNoteDTEExtension{
		NombreEntrega:    extension.GetDeliveryName(),
		DocumentoEntrega: extension.GetDeliveryDocument(),
		NombreRecibe:     extension.GetReceiverName(),
		DocumentoRecibe:  extension.GetReceiverDocument(),
		Observacion:      extension.GetObservation(),
	}
}
