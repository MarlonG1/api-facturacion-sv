package response_mapper

import (
	commonModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/invalidation"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToMHInvalidation(doc *models.InvalidationDocument) *structs.InvalidationResponse {
	if doc == nil {
		return nil
	}

	return &structs.InvalidationResponse{
		Identificacion: *MapIdentificationResponse(doc.Identification),
		Emisor:         *MapIssuerResponse(doc.Issuer),
		Documento:      *invalidation.MapInvalidatedDocumentResponse(doc.Document),
		Motivo:         *invalidation.MapInvalidationReasonResponse(doc.Reason),
	}
}

func MapIdentificationResponse(identification *commonModels.Identification) *structs.InvalidationIdentification {
	if identification == nil {
		return nil
	}

	return &structs.InvalidationIdentification{
		Version:          identification.Version.GetValue(),
		Ambiente:         identification.Ambient.GetValue(),
		CodigoGeneracion: identification.GenerationCode.GetValue(),
		FecAnula:         identification.EmissionDate.GetValue().Format("2006-01-02"),
		HorAnula:         identification.EmissionTime.GetValue().Format("15:04:05"),
	}
}

func MapIssuerResponse(issuer *commonModels.Issuer) *structs.InvalidationIssuer {
	if issuer == nil {
		return nil
	}

	return &structs.InvalidationIssuer{
		NIT:                   issuer.NIT.GetValue(),
		Nombre:                issuer.Name,
		TipoEstablecimiento:   issuer.EstablishmentType.GetValue(),
		NombreComercial:       &issuer.CommercialName,
		Telefono:              issuer.Phone.GetValue(),
		Correo:                issuer.Email.GetValue(),
		CodigoEstablecimiento: issuer.EstablishmentCode,
		POSCodigo:             issuer.POSCode,
	}
}
