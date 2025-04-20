package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// CreateDefaultAddress crea una dirección predeterminada válida
func CreateDefaultAddress() *structs.AddressRequest {
	return &structs.AddressRequest{
		Department:   "06",
		Municipality: "20",
		Complement:   "Colonia Escalón, Calle La Reforma #123, San Salvador",
	}
}

// CreateDefaultAppendix crea un apéndice predeterminado válido
func CreateDefaultAppendix() structs.AppendixRequest {
	return structs.AppendixRequest{
		Field: "nota_interna",
		Label: "Nota interna",
		Value: "Información adicional para el documento",
	}
}

// CreateDefaultExtension crea una extensión predeterminada válida
func CreateDefaultExtension() *structs.ExtensionRequest {
	observation := "Observación de prueba"
	vehiculePlate := "P123-456"

	return &structs.ExtensionRequest{
		DeliveryName:     "Juan Pérez",
		DeliveryDocument: "12345678-9",
		ReceiverName:     "Ana López",
		ReceiverDocument: "98765432-1",
		Observation:      &observation,
		VehiculePlate:    &vehiculePlate,
	}
}

// CreateDefaultReceiver crea un receptor predeterminado válido
func CreateDefaultReceiver() *structs.ReceiverRequest {
	docType := "36" // NIT
	docNumber := "06141804941035"
	name := "Empresa Servicios Generales, S.A. de C.V."
	nrc := "123456"
	phone := "22123456"
	email := "empresa@example.com"
	activityCode := "46900"
	activityDesc := "Venta al por mayor de otros productos"
	commercialName := "ServiGeneral"

	return &structs.ReceiverRequest{
		DocumentType:   &docType,
		DocumentNumber: &docNumber,
		Name:           &name,
		NRC:            &nrc,
		Address:        CreateDefaultAddress(),
		Phone:          &phone,
		Email:          &email,
		ActivityCode:   &activityCode,
		ActivityDesc:   &activityDesc,
		CommercialName: &commercialName,
	}
}

// CreateDefaultReceiverWithoutDocsFields crea un receptor predeterminado válido sin campos de tipo documento y numero documento
func CreateDefaultReceiverWithoutDocsFields() *structs.ReceiverRequest {
	name := "Empresa Servicios Generales, S.A. de C.V."
	nrc := "123456"
	nit := "06141804941035"
	phone := "22123456"
	email := "empresa@example.com"
	activityCode := "46900"
	activityDesc := "Venta al por mayor de otros productos"
	commercialName := "ServiGeneral"

	return &structs.ReceiverRequest{
		Name:           &name,
		NRC:            &nrc,
		NIT:            &nit,
		Address:        CreateDefaultAddress(),
		Phone:          &phone,
		Email:          &email,
		ActivityCode:   &activityCode,
		ActivityDesc:   &activityDesc,
		CommercialName: &commercialName,
	}
}

// CreateDefaultPayment crea un pago predeterminado válido
func CreateDefaultPayment() structs.PaymentRequest {
	reference := "REF-123"

	return structs.PaymentRequest{
		Code:      "01",
		Amount:    100.0,
		Reference: &reference,
	}
}

// CreateDefaultThirdPartySale crea una venta a terceros predeterminada válida
func CreateDefaultThirdPartySale() *structs.ThirdPartySaleRequest {
	return &structs.ThirdPartySaleRequest{
		NIT:  "06141804941035",
		Name: "Empresa Tercero S.A. de C.V.",
	}
}

// CreateDefaultRelatedDocument crea un documento relacionado predeterminado válido
func CreateDefaultRelatedDocument() structs.RelatedDocRequest {
	return structs.RelatedDocRequest{
		DocumentType:   "03", // CCF
		GenerationType: 1,    // Normal
		DocumentNumber: "S221001346",
		EmissionDate:   "2025-03-22",
	}
}

// CreateDefaultIssuer crea un emisor por defecto para uso en pruebas
func CreateDefaultIssuer() *dte.IssuerDTE {
	return &dte.IssuerDTE{
		NIT:                  "11111111111111",
		NRC:                  "1111111",
		CommercialName:       "EJEMPLO SA",
		BusinessName:         "EMPRESA DE PRUEBAS SA DE CV",
		EconomicActivity:     "11111",
		EconomicActivityDesc: "Venta al por mayor de otros productos",
		EstablishmentCode:    utils.ToStringPointer("C001"),
		Email:                utils.ToStringPointer("email@gmail.com"),
		Phone:                utils.ToStringPointer("22567890"),
		Address: &user.Address{
			Department:   "06",
			Municipality: "20",
			Complement:   "BOULEVARD SANTA ELENA SUR, SANTA TECLA",
		},
		EstablishmentType:   "02",
		EstablishmentCodeMH: nil,
		POSCode:             nil,
		POSCodeMH:           nil,
	}
}

// CreateCustomIssuer crea un emisor personalizado para uso en pruebas
func CreateCustomIssuer(nit, nrc, businessName string) *dte.IssuerDTE {
	issuer := CreateDefaultIssuer()
	issuer.NIT = nit
	issuer.NRC = nrc
	issuer.BusinessName = businessName
	return issuer
}

// CreateDefaultOtherDocument crea un documento adicional predeterminado válido
func CreateDefaultOtherDocument() structs.OtherDocRequest {
	description := "Documento adicional"
	detail := "Detalle del documento adicional"

	return structs.OtherDocRequest{
		DocumentCode: 1,
		Description:  &description,
		Detail:       &detail,
	}
}

// CreateIdentification crea una identificación con tipo y versión específicos
func CreateIdentification(dteType string, version int) (*models.Identification, error) {
	now := utils.TimeNow()

	versionValue := document.NewValidatedVersion(version)
	dteTypeValue := document.NewValidatedDTEType(dteType)
	currency := financial.NewValidatedCurrency("USD")

	ambient, err := document.NewAmbient()
	if err != nil {
		return nil, err
	}

	emissionDate, err := temporal.NewEmissionDate(now)
	if err != nil {
		return nil, err
	}

	emissionTime, err := temporal.NewEmissionTime(now)
	if err != nil {
		return nil, err
	}

	modelType, err := document.NewModelType(1) // Modelo normal
	if err != nil {
		return nil, err
	}

	operationType, err := document.NewOperationType(1) // Transmisión normal
	if err != nil {
		return nil, err
	}

	return &models.Identification{
		Version:       *versionValue,
		Ambient:       *ambient,
		DTEType:       *dteTypeValue,
		Currency:      *currency,
		OperationType: *operationType,
		ModelType:     *modelType,
		EmissionDate:  *emissionDate,
		EmissionTime:  *emissionTime,
	}, nil
}

// CreateIdentificationWithInvalidVersion crea una identificación con versión inválida
func CreateIdentificationWithInvalidVersion(dteType string) (*models.Identification, error) {
	id, err := CreateIdentification(dteType, 1)
	if err != nil {
		return nil, err
	}

	invalidVersion := document.NewValidatedVersion(99) // Versión inválida
	id.Version = *invalidVersion

	return id, nil
}

// CreateIdentificationWithInvalidDTEType crea una identificación con tipo DTE inválido
func CreateIdentificationWithInvalidDTEType() (*models.Identification, error) {
	id, err := CreateIdentification("01", 1)
	if err != nil {
		return nil, err
	}

	invalidType := document.NewValidatedDTEType("99") // Tipo DTE inválido
	id.DTEType = *invalidType

	return id, nil
}
