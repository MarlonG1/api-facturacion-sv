package fixtures

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/invalidation_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// InvalidationBuilder - Builder especializado para documentos de invalidación
type InvalidationBuilder struct {
	document *invalidation_models.InvalidationDocument
	err      error
}

func NewInvalidationBuilder() *InvalidationBuilder {
	return &InvalidationBuilder{
		document: &invalidation_models.InvalidationDocument{
			Identification: &models.Identification{},
			Issuer:         &models.Issuer{},
			Document:       &invalidation_models.InvalidatedDocument{},
			Reason:         &invalidation_models.InvalidationReason{},
		},
		err: nil,
	}
}

func (b *InvalidationBuilder) Document() *invalidation_models.InvalidationDocument {
	return b.document
}

func (b *InvalidationBuilder) Build() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.document, nil
}

func (b *InvalidationBuilder) BuildWithoutValidation() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.document, nil
}

func (b *InvalidationBuilder) setError(err error) *InvalidationBuilder {
	if b.err == nil && err != nil {
		b.err = err
	}
	return b
}

func (b *InvalidationBuilder) AddIdentification() *InvalidationBuilder {
	if b.err != nil {
		return b
	}

	// Crear identificación básica
	b.setError(b.document.Identification.SetVersion(1))
	b.setError(b.document.Identification.SetAmbient(constants.Testing))
	b.setError(b.document.Identification.SetDTEType(constants.FacturaElectronica))
	b.setError(b.document.Identification.SetControlNumber("DTE-05-00000001-000000000000001"))
	b.setError(b.document.Identification.GenerateCode())
	b.setError(b.document.Identification.SetModelType(constants.ModeloFacturacionPrevio))
	b.setError(b.document.Identification.SetOperationType(1))
	b.setError(b.document.Identification.SetEmissionDate(utils.TimeNow()))
	b.setError(b.document.Identification.SetEmissionTime(utils.TimeNow()))
	b.setError(b.document.Identification.SetCurrency("USD"))

	return b
}

func (b *InvalidationBuilder) AddIssuer() *InvalidationBuilder {
	if b.err != nil {
		return b
	}

	// Crear datos del emisor
	b.setError(b.document.Issuer.SetNIT("12345678901234"))
	b.setError(b.document.Issuer.SetNRC("12345678"))
	b.setError(b.document.Issuer.SetName("EMPRESA EMISORA, S.A. DE C.V."))
	b.setError(b.document.Issuer.SetActivityCode("12345"))
	b.setError(b.document.Issuer.SetActivityDescription("Venta de productos electrónicos"))
	b.setError(b.document.Issuer.SetEstablishmentType(constants.CasaMatriz))

	// Crear y configurar la dirección
	address := &models.Address{}
	b.setError(address.SetDepartment("06"))
	b.setError(address.SetMunicipality("21"))
	b.setError(address.SetComplement("Calle Principal, Edificio Central #123"))

	// Asignar dirección al emisor
	if b.err == nil {
		b.setError(b.document.Issuer.SetAddress(address))
	}

	b.setError(b.document.Issuer.SetPhone("22225555"))
	b.setError(b.document.Issuer.SetEmail("info@google.com"))
	b.setError(b.document.Issuer.SetCommercialName("EMPRESA TECH"))

	// Establecer campos opcionales
	establishmentCode := "001"
	establishmentMHCode := "EST001"
	posCode := "POS01"
	posMHCode := "POS001"
	b.setError(b.document.Issuer.SetEstablishmentCode(&establishmentCode))
	b.setError(b.document.Issuer.SetEstablishmentMHCode(&establishmentMHCode))
	b.setError(b.document.Issuer.SetPOSCode(&posCode))
	b.setError(b.document.Issuer.SetPOSMHCode(&posMHCode))

	return b
}

func (b *InvalidationBuilder) AddInvalidatedDocument() *InvalidationBuilder {
	if b.err != nil {
		return b
	}

	// Establecer tipo de documento a invalidar
	docType, err := document.NewDTEType(constants.FacturaElectronica)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.Type = *docType

	// Establecer código de generación
	generationCode, err := identification.NewGenerationCode()
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.GenerationCode = *generationCode

	// Establecer número de control
	controlNumber, err := identification.NewControlNumber("DTE-01-00000000-000000000000001")
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.ControlNumber = *controlNumber

	// Establecer sello de recepción
	b.document.Document.ReceptionStamp = "2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM"

	// Establecer fecha de emisión
	emissionDate, err := temporal.NewEmissionDate(time.Now().Add(-24 * time.Hour))
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.EmissionDate = *emissionDate

	// Establecer monto de IVA
	ivaAmount, err := financial.NewAmount(13.00)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.IVAAmount = ivaAmount

	// Datos del receptor del documento original
	docTypeReceiver, err := document.NewDTEType(constants.FacturaElectronica)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.DocumentType = docTypeReceiver

	docNumber, err := identification.NewDocumentNumber("01234567-8", constants.DUI)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.DocumentNumber = docNumber

	name := "CLIENTE EJEMPLO"
	b.document.Document.Name = &name

	email, err := base.NewEmail("cliente@ejemplo.com")
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.Email = email

	phone, err := base.NewPhone("77778888")
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.Phone = phone

	return b
}

func (b *InvalidationBuilder) AddReplacementCode() *InvalidationBuilder {
	if b.err != nil {
		return b
	}

	// Solo necesario para invalidaciones tipo 1 (reemplazo)
	replacementCode, err := identification.NewGenerationCode()
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Document.ReplacementCode = replacementCode

	return b
}

func (b *InvalidationBuilder) AddInvalidationReason(invalidationType int) *InvalidationBuilder {
	if b.err != nil {
		return b
	}

	// Establecer tipo de invalidación
	invalidationTypeObj, err := document.NewInvalidationType(invalidationType)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Reason.Type = *invalidationTypeObj

	// Datos del responsable
	b.document.Reason.ResponsibleName = "JUAN RESPONSABLE"

	responsibleDocType, err := document.NewDTETypeForReceiver(constants.DUI)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Reason.ResponsibleDocType = *responsibleDocType

	responsibleDocNum, err := identification.NewDocumentNumber("01234567-8", constants.DUI)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Reason.ResponsibleDocNum = *responsibleDocNum

	// Datos del solicitante
	b.document.Reason.RequesterName = "ANA SOLICITANTE"

	requesterDocType, err := document.NewDTETypeForReceiver(constants.DUI)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Reason.RequesterDocType = *requesterDocType

	requesterDocNum, err := identification.NewDocumentNumber("98765432-1", constants.DUI)
	if err != nil {
		b.setError(err)
		return b
	}
	b.document.Reason.RequesterDocNum = *requesterDocNum

	// Si es tipo 3 (definitiva), agregar razón
	if invalidationType == 3 {
		reasonText := "Documento con errores graves que impiden su utilización"
		invalidReason, err := document.NewInvalidationReason(reasonText)
		if err != nil {
			b.setError(err)
			return b
		}
		b.document.Reason.Reason = invalidReason
	} else {
		// Para tipos 1 y 2, la razón debe ser nil
		b.document.Reason.Reason = nil
	}

	return b
}

// BuildInvalidationWithReplacement construye un documento de invalidación tipo 1 (con reemplazo)
func BuildInvalidationWithReplacement() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddReplacementCode().
		AddInvalidationReason(1)

	return builder.Build()
}

// BuildInvalidationWithAnnulment construye un documento de invalidación tipo 2 (anulación)
func BuildInvalidationWithAnnulment() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddInvalidationReason(2)

	return builder.Build()
}

// BuildInvalidationDefinitive construye un documento de invalidación tipo 3 (definitiva)
func BuildInvalidationDefinitive() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddReplacementCode(). // Algunas invalidaciones tipo 3 pueden tener código de reemplazo
		AddInvalidationReason(3)

	return builder.Build()
}

// BuildInvalidInvalidation construye un documento de invalidación inválido (tipo 2 con código de reemplazo)
func BuildInvalidInvalidation() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddReplacementCode(). // Inválido para tipo 2
		AddInvalidationReason(2)

	return builder.BuildWithoutValidation()
}

// BuildInvalidationWithInvalidReason construye un documento de invalidación inválido por tener razón para tipo 1/2
func BuildInvalidationWithInvalidReason() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddInvalidationReason(1)

	// Agregar razón incorrectamente para tipo 1
	reasonText := "Esta razón no debería existir para tipo 1"
	invalidReason, _ := document.NewInvalidationReason(reasonText)
	builder.document.Reason.Reason = invalidReason

	return builder.BuildWithoutValidation()
}

// BuildInvalidationWithMissingReason construye una invalidación tipo 3 sin razón (inválido)
func BuildInvalidationWithMissingReason() (*invalidation_models.InvalidationDocument, error) {
	builder := NewInvalidationBuilder()

	builder.AddIdentification().
		AddIssuer().
		AddInvalidatedDocument().
		AddReplacementCode()

	// Crear tipo 3 pero sin especificar razón
	invalidationType, _ := document.NewInvalidationType(3)
	builder.document.Reason.Type = *invalidationType

	// Datos del responsable
	builder.document.Reason.ResponsibleName = "JUAN RESPONSABLE"

	responsibleDocType, _ := document.NewDTETypeForReceiver(constants.DUI)
	builder.document.Reason.ResponsibleDocType = *responsibleDocType

	responsibleDocNum, _ := identification.NewDocumentNumber("01234567-8", constants.DUI)
	builder.document.Reason.ResponsibleDocNum = *responsibleDocNum

	// Datos del solicitante
	builder.document.Reason.RequesterName = "ANA SOLICITANTE"

	requesterDocType, _ := document.NewDTETypeForReceiver(constants.DUI)
	builder.document.Reason.RequesterDocType = *requesterDocType

	requesterDocNum, _ := identification.NewDocumentNumber("98765432-1", constants.DUI)
	builder.document.Reason.RequesterDocNum = *requesterDocNum

	// Dejar razón como nil (inválido para tipo 3)
	builder.document.Reason.Reason = nil

	return builder.BuildWithoutValidation()
}

// BuildInvalidation es un método genérico para construir una invalidación válida
func BuildInvalidation() (*invalidation_models.InvalidationDocument, error) {
	return BuildInvalidationWithReplacement()
}
