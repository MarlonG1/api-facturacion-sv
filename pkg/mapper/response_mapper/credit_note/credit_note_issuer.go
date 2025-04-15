package credit_note

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCreditNoteIssuer mapea el emisor de una invoice electrónica a un modelo de emisor -> Origen: Response
func MapCreditNoteIssuer(issuer interfaces.Issuer) structs.CreditNoteDTEIssuer {
	result := structs.CreditNoteDTEIssuer{
		NIT:                 issuer.GetNIT(),
		NRC:                 issuer.GetNRC(),
		Nombre:              issuer.GetName(),
		CodActividad:        issuer.GetActivityCode(),
		DescActividad:       issuer.GetActivityDescription(),
		TipoEstablecimiento: issuer.GetEstablishmentType(),
		Direccion:           common.MapCommonResponseAddress(issuer.GetAddress()),
		Telefono:            issuer.GetPhone(),
		Correo:              issuer.GetEmail(),
	}

	// Mapear campos opcionales si tienen valor
	if name := issuer.GetCommercialName(); name != "" {
		result.NombreComercial = &name
	}

	return result
}
