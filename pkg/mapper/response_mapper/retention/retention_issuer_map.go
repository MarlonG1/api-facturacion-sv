package retention

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func MapRetentionResponseIssuer(issuer interfaces.Issuer) structs.RetentionIssuer {
	result := structs.RetentionIssuer{
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
	if code := issuer.GetEstablishmentCode(); code != nil {
		result.CodEstable = code
	}
	if code := issuer.GetEstablishmentMHCode(); code != nil {
		result.CodEstableMH = code
	}
	if code := issuer.GetPOSCode(); code != nil {
		result.CodPuntoVenta = code
	}
	if code := issuer.GetPOSMHCode(); code != nil {
		result.CodPuntoVentaMH = code
	}

	return result
}
