package fixtures

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invalidation/invalidation_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type DTEBuilder struct {
	document *models.DTEDocument
	err      error
}

// NewDTEBuilder crea un nuevo builder para DTEDocument
func NewDTEBuilder() *DTEBuilder {
	return &DTEBuilder{
		document: &models.DTEDocument{
			Items:            make([]interfaces.Item, 0),
			Appendix:         make([]interfaces.Appendix, 0),
			RelatedDocuments: make([]interfaces.RelatedDocument, 0),
			OtherDocuments:   make([]interfaces.OtherDocuments, 0),
		},
	}
}

// Build construye y valida el documento DTE, devolviendo cualquier error que haya ocurrido
func (b *DTEBuilder) Build() (*models.DTEDocument, error) {
	// Si ya hubo un error durante la construcción, devolverlo
	if b.err != nil {
		return nil, b.err
	}

	// Realizar validaciones completas
	err := b.document.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := b.document.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return b.document, nil
}

// BuildWithoutValidation construye el documento sin validaciones
func (b *DTEBuilder) BuildWithoutValidation() (*models.DTEDocument, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.document, nil
}

// setError es un método auxiliar que establece el error si aún no se ha establecido
func (b *DTEBuilder) setError(err error) *DTEBuilder {
	if b.err == nil && err != nil {
		b.err = err
	}
	return b
}

// AddIdentification añade datos de identificación predeterminados válidos
func (b *DTEBuilder) AddIdentification() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear un objeto Identification
	identification := &models.Identification{}

	// Establecer los valores usando los setters
	b.setError(identification.SetVersion(1))
	b.setError(identification.SetAmbient(constants.Testing))
	b.setError(identification.SetDTEType(constants.FacturaElectronica))
	b.setError(identification.SetControlNumber("DTE-01-12345678-123456789012345"))
	b.setError(identification.GenerateCode())
	b.setError(identification.SetModelType(constants.ModeloFacturacionPrevio))
	b.setError(identification.SetOperationType(1))
	b.setError(identification.SetEmissionDate(utils.TimeNow()))
	b.setError(identification.SetEmissionTime(utils.TimeNow()))
	b.setError(identification.SetCurrency("USD"))

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetIdentification(identification))
	}

	return b
}

// AddIdentificationWithContingency añade datos de identificación con contingencia
func (b *DTEBuilder) AddIdentificationWithContingency() *DTEBuilder {
	// Primero añadir la identificación base
	b.AddIdentification()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener la identificación y modificarla para contingencia
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if !ok || identification == nil {
		b.setError(fmt.Errorf("failed to get identification"))
		return b
	}

	// Modificar para contingencia usando setters
	b.setError(identification.SetOperationType(constants.TransmisionContingencia))
	b.setError(identification.SetModelType(constants.ModeloFacturacionDiferido))

	// Establecer tipo de contingencia y razón
	contingencyType := constants.FallaServicioInternet
	contingencyReason := "Falla de servicio de internet del proveedor"
	b.setError(identification.SetContingencyType(&contingencyType))
	b.setError(identification.SetContingencyReason(&contingencyReason))

	return b
}

// AddIssuer añade datos de emisor predeterminados válidos
func (b *DTEBuilder) AddIssuer() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear un nuevo emisor
	issuer := &models.Issuer{}

	// Establecer los valores usando los setters
	b.setError(issuer.SetNIT("12345678901234"))
	b.setError(issuer.SetNRC("12345678"))
	b.setError(issuer.SetName("COMPANY EXAMPLE, S.A. DE C.V."))
	b.setError(issuer.SetActivityCode("12345"))
	b.setError(issuer.SetActivityDescription("Electronic products sales"))
	b.setError(issuer.SetEstablishmentType(constants.CasaMatriz))

	// Crear y configurar la dirección
	address := &models.Address{}
	b.setError(address.SetDepartment("06"))
	b.setError(address.SetMunicipality("21"))
	b.setError(address.SetComplement("Example Street, Central Building #123"))

	// Asignar la dirección al emisor
	if b.err == nil {
		b.setError(issuer.SetAddress(address))
	}

	b.setError(issuer.SetPhone("22225555"))
	b.setError(issuer.SetEmail("info@google.com"))
	b.setError(issuer.SetCommercialName("ELECTRO STORE"))

	// Establecer campos opcionales
	establishmentCode := "001"
	establishmentMHCode := "EST001"
	posCode := "POS01"
	posMHCode := "POS001"
	b.setError(issuer.SetEstablishmentCode(&establishmentCode))
	b.setError(issuer.SetEstablishmentMHCode(&establishmentMHCode))
	b.setError(issuer.SetPOSCode(&posCode))
	b.setError(issuer.SetPOSMHCode(&posMHCode))

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetIssuer(issuer))
	}

	return b
}

// AddIssuerWithInvalidNIT añade datos de emisor con NIT inválido para testing
func (b *DTEBuilder) AddIssuerWithInvalidNIT() *DTEBuilder {
	// Primero añadir el emisor base
	b.AddIssuer()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener el emisor y modificar el NIT
	issuer, ok := b.document.GetIssuer().(*models.Issuer)
	if !ok || issuer == nil {
		b.setError(fmt.Errorf("failed to get issuer"))
		return b
	}

	// Intentar establecer un NIT inválido (menos dígitos)
	b.setError(issuer.SetNIT("123456789"))

	return b
}

// AddReceiver añade datos de receptor predeterminados válidos para factura
func (b *DTEBuilder) AddReceiver() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear un nuevo receptor
	receiver := &models.Receiver{}

	// Establecer valores
	name := "John Albert Smith"
	email := "john.smith@email.com"
	phone := "77778888"
	docType := constants.DUI
	docNumber := "01234567-8"

	b.setError(receiver.SetName(&name))
	b.setError(receiver.SetDocumentType(&docType))
	b.setError(receiver.SetDocumentNumber(&docNumber))
	b.setError(receiver.SetEmail(&email))
	b.setError(receiver.SetPhone(&phone))

	// Crear y configurar la dirección
	address := &models.Address{}
	b.setError(address.SetDepartment("06"))
	b.setError(address.SetMunicipality("22"))
	b.setError(address.SetComplement("Example Neighborhood, House #456"))

	// Asignar la dirección al receptor
	if b.err == nil {
		b.setError(receiver.SetAddress(address))
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetReceiver(receiver))
	}

	return b
}

// AddReceiverForCCF añade datos de receptor predeterminados válidos para CCF
func (b *DTEBuilder) AddReceiverForCCF() *DTEBuilder {
	// Primero añadir el receptor base
	b.AddReceiver()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener el receptor y modificarlo para CCF
	receiver, ok := b.document.GetReceiver().(*models.Receiver)
	if !ok || receiver == nil {
		b.setError(fmt.Errorf("failed to get receiver"))
		return b
	}

	// Agregar información específica para CCF
	nrc := "987654"
	activityDescription := "Purchase of goods and services"
	activityCode := "56789"
	commercialName := "CLIENT COMPANY INC."
	nit := "98765432101234"

	b.setError(receiver.SetNRC(&nrc))
	b.setError(receiver.SetActivityDescription(&activityDescription))
	b.setError(receiver.SetActivityCode(&activityCode))
	b.setError(receiver.SetCommercialName(&commercialName))
	b.setError(receiver.SetNIT(&nit))

	return b
}

// AddReceiverWithNoNRC añade datos de receptor sin NRC (inválido para CCF)
func (b *DTEBuilder) AddReceiverWithNoNRC() *DTEBuilder {
	// Añadir el receptor base que no tiene NRC
	b.AddReceiver()
	return b
}

// AddItems añade items predeterminados válidos
func (b *DTEBuilder) AddItems() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear items
	items := make([]interfaces.Item, 0, 2)

	// Primer item (producto)
	item1 := &models.Item{}
	b.setError(item1.SetNumber(1))
	b.setError(item1.SetType(constants.Producto))
	b.setError(item1.SetDescription("HP EliteBook 840 G8 Laptop"))
	b.setError(item1.SetQuantity(1.0))
	b.setError(item1.SetUnitMeasure(1))
	b.setError(item1.SetUnitPrice(899.99))
	b.setError(item1.SetDiscount(0.0))
	b.setError(item1.SetTaxes([]string{constants.TaxIVA}))
	b.setError(item1.SetItemCode("LAPTOP-HP-001"))

	if b.err == nil {
		items = append(items, item1)
	}

	// Segundo item (servicio)
	item2 := &models.Item{}
	b.setError(item2.SetNumber(2))
	b.setError(item2.SetType(constants.Servicio))
	b.setError(item2.SetDescription("Software configuration and installation"))
	b.setError(item2.SetQuantity(1.0))
	b.setError(item2.SetUnitMeasure(1))
	b.setError(item2.SetUnitPrice(50.00))
	b.setError(item2.SetDiscount(0.0))
	b.setError(item2.SetTaxes([]string{constants.TaxIVA}))
	b.setError(item2.SetItemCode("Licencia de Windows 11"))

	if b.err == nil {
		items = append(items, item2)
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetItems(items))
	}

	return b
}

// AddItemsWithDiscount añade items con descuentos
func (b *DTEBuilder) AddItemsWithDiscount() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear items
	items := make([]interfaces.Item, 0, 2)

	// Primer item (producto con descuento)
	item1 := &models.Item{}
	b.setError(item1.SetNumber(1))
	b.setError(item1.SetType(constants.Producto))
	b.setError(item1.SetDescription("HP EliteBook 840 G8 Laptop"))
	b.setError(item1.SetQuantity(1.0))
	b.setError(item1.SetUnitMeasure(1))
	b.setError(item1.SetUnitPrice(899.99))
	b.setError(item1.SetDiscount(10.0)) // 10% de descuento
	b.setError(item1.SetTaxes([]string{constants.TaxIVA}))
	b.setError(item1.SetItemCode("LAPTOP-HP-001"))

	if b.err == nil {
		items = append(items, item1)
	}

	// Segundo item (servicio con descuento)
	item2 := &models.Item{}
	b.setError(item2.SetNumber(2))
	b.setError(item2.SetType(constants.Servicio))
	b.setError(item2.SetDescription("Software configuration and installation"))
	b.setError(item2.SetQuantity(1.0))
	b.setError(item2.SetUnitMeasure(1))
	b.setError(item2.SetUnitPrice(50.00))
	b.setError(item2.SetDiscount(5.0)) // 5% de descuento
	b.setError(item2.SetTaxes([]string{constants.TaxIVA}))

	if b.err == nil {
		items = append(items, item2)
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetItems(items))
	}

	return b
}

// AddItemsWithInvalidTax añade items con impuestos inválidos para testing
func (b *DTEBuilder) AddItemsWithInvalidTax() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear item con código de impuesto inválido
	invalidItem := &models.Item{}
	b.setError(invalidItem.SetNumber(1))
	b.setError(invalidItem.SetType(constants.Producto))
	b.setError(invalidItem.SetDescription("Product with invalid tax"))
	b.setError(invalidItem.SetQuantity(1.0))
	b.setError(invalidItem.SetUnitMeasure(1))
	b.setError(invalidItem.SetUnitPrice(100.00))
	b.setError(invalidItem.SetDiscount(0.0))
	b.setError(invalidItem.SetTaxes([]string{"ZZ"})) // Código inválido

	// Asignar al documento
	if b.err == nil {
		items := []interfaces.Item{invalidItem}
		b.setError(b.document.SetItems(items))
	}

	return b
}

// AddSummary añade resumen predeterminado válido
func (b *DTEBuilder) AddSummary() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Calcular valores basados en los items
	var totalTaxed float64
	for _, i := range b.document.GetItems() {
		// Extraer el ítem
		item, ok := i.(*models.Item)
		if !ok {
			continue
		}

		// Calcular precio con descuento
		discountFactor := 1.0 - (item.GetDiscount() / 100.0)
		itemTotal := item.GetQuantity() * item.GetUnitPrice() * discountFactor
		totalTaxed += itemTotal
	}

	// Calcular IVA y total
	iva := totalTaxed * 0.13
	totalOperation := totalTaxed + iva

	// Crear el resumen
	summary := &models.Summary{}

	// Establecer valores
	b.setError(summary.SetTotalNonSubject(0.00))
	b.setError(summary.SetTotalExempt(0.00))
	b.setError(summary.SetTotalTaxed(totalTaxed))
	b.setError(summary.SetSubTotal(totalTaxed))
	b.setError(summary.SetSubtotalSales(totalTaxed))
	b.setError(summary.SetNonSubjectDiscount(0.00))
	b.setError(summary.SetExemptDiscount(0.00))
	b.setError(summary.SetDiscountPercentage(0.00))
	b.setError(summary.SetTotalDiscount(0.00))
	b.setError(summary.SetTotalOperation(totalOperation))
	b.setError(summary.SetTotalNotTaxed(0.00))
	b.setError(summary.SetOperationCondition(constants.Cash))
	b.setError(summary.SetTotalToPay(totalOperation))
	b.setError(summary.SetTotalInWords("UN MIL DOLARES CON 00/100 CENTAVOS"))

	// Crear impuesto IVA
	tax := &models.Tax{}
	b.setError(tax.SetCode(constants.TaxIVA))
	b.setError(tax.SetDescription("IVA 13%"))
	b.setError(tax.SetValue(iva))

	// Añadir impuesto al resumen
	var taxes []interfaces.Tax
	if b.err == nil {
		taxes = append(taxes, tax)
		b.setError(summary.SetTotalTaxes(taxes))
	}

	// Crear tipo de pago en efectivo
	payment := &models.PaymentType{}
	b.setError(payment.SetCode(constants.BilletesMonedas))
	b.setError(payment.SetAmount(totalOperation))
	b.setError(payment.SetReference("Cash payment"))

	// Añadir tipo de pago al resumen
	var payments []interfaces.PaymentType
	if b.err == nil {
		payments = append(payments, payment)
		b.setError(summary.SetPaymentTypes(payments))
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetSummary(summary))
	}

	return b
}

// AddSummaryWithCredit añade resumen con condición de pago a crédito
func (b *DTEBuilder) AddSummaryWithCredit() *DTEBuilder {
	// Primero añadir el resumen base
	b.AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener el resumen
	summary, ok := b.document.GetSummary().(*models.Summary)
	if !ok || summary == nil {
		b.setError(fmt.Errorf("failed to get summary"))
		return b
	}

	// Modificar la condición de operación a crédito
	b.setError(summary.SetOperationCondition(constants.Credit))

	// Modificar el tipo de pago
	if len(summary.GetPaymentTypes()) > 0 {
		paymentType, ok := summary.GetPaymentTypes()[0].(*models.PaymentType)
		if ok {
			// Cambiar a transferencia bancaria
			b.setError(paymentType.SetCode(constants.TransBancaria))
			b.setError(paymentType.SetReference("Bank transfer"))

			// Agregar términos de pago
			term := "01" // 30 días
			period := 30
			b.setError(paymentType.SetTerm(&term))
			b.setError(paymentType.SetPeriod(&period))
		}
	}

	return b
}

// AddSummaryWithInvalidPayment añade resumen con pagos inválidos para testing
func (b *DTEBuilder) AddSummaryWithInvalidPayment() *DTEBuilder {
	// Primero añadir el resumen base
	b.AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener el resumen
	summary, ok := b.document.GetSummary().(*models.Summary)
	if !ok || summary == nil {
		b.setError(fmt.Errorf("failed to get summary"))
		return b
	}

	// Modificar la condición a crédito pero dejar el pago en efectivo (inválido)
	b.setError(summary.SetOperationCondition(constants.Credit))

	return b
}

// AddSummaryWithInvalidTotal añade resumen con total incorrecto para testing
func (b *DTEBuilder) AddSummaryWithInvalidTotal() *DTEBuilder {
	// Primero añadir el resumen base
	b.AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Obtener el resumen
	summary, ok := b.document.GetSummary().(*models.Summary)
	if !ok || summary == nil {
		b.setError(fmt.Errorf("failed to get summary"))
		return b
	}

	// Modificar el total para que no coincida con la suma de los items
	b.setError(summary.SetTotalOperation(1000.00))
	b.setError(summary.SetTotalToPay(1000.00))

	// Actualizar los pagos para que coincidan con el nuevo total (para evitar error de pago != total)
	if len(summary.GetPaymentTypes()) > 0 {
		paymentType, ok := summary.GetPaymentTypes()[0].(*models.PaymentType)
		if ok {
			b.setError(paymentType.SetAmount(1000.00))
		}
	}

	return b
}

// AddExtension añade extensión predeterminada válida
func (b *DTEBuilder) AddExtension() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear extensión
	extension := &models.Extension{}

	// Establecer valores
	b.setError(extension.SetDeliveryName("Mary Rodriguez"))
	b.setError(extension.SetDeliveryDocument("12345678-9"))
	b.setError(extension.SetReceiverName("Louis Gonzalez"))
	b.setError(extension.SetReceiverDocument("98765432-1"))

	vehiculePlate := "P123-456"
	observation := "Delivery at building reception. Contact the recipient."
	b.setError(extension.SetVehiculePlate(&vehiculePlate))
	b.setError(extension.SetObservation(&observation))

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetExtension(extension))
	}

	return b
}

// AddAppendixes añade apéndices predeterminados válidos
func (b *DTEBuilder) AddAppendixes() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	appendixCount := 3
	appendixesInterfaces := make([]interfaces.Appendix, 0, appendixCount)

	for i := 0; i < appendixCount; i++ {
		appendix := &models.Appendix{}

		// Establecer valores
		b.setError(appendix.SetField(fmt.Sprintf("NOTE%d", i+1)))
		b.setError(appendix.SetLabel(fmt.Sprintf("Additional Information %d", i+1)))
		b.setError(appendix.SetValue(fmt.Sprintf("Additional information content %d", i+1)))

		if b.err == nil {
			appendixesInterfaces = append(appendixesInterfaces, appendix)
		}
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetAppendix(appendixesInterfaces))
	}

	return b
}

// AddRelatedDocuments añade documentos relacionados predeterminados válidos
func (b *DTEBuilder) AddRelatedDocuments() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear documento relacionado
	relatedDoc := &models.RelatedDocument{}

	// Establecer valores
	b.setError(relatedDoc.SetDocumentType(constants.NotaRemisionElectronica))
	b.setError(relatedDoc.SetGenerationType(constants.ElectronicDocument))
	b.setError(relatedDoc.SetDocumentNumber("DA1E261A-BAD7-460F-AD15-04F2E281FC6A"))
	b.setError(relatedDoc.SetEmissionDate(utils.TimeNow().Add(-24 * time.Hour))) // Ayer

	// Añadir al documento
	if b.err == nil {
		relatedDocs := make([]interfaces.RelatedDocument, 0, 1)
		relatedDocs = append(relatedDocs, relatedDoc)
		b.setError(b.document.SetRelatedDocuments(relatedDocs))
	}

	return b
}

// AddInvalidRelatedDocument añade documentos relacionados inválidos para testing
func (b *DTEBuilder) AddInvalidRelatedDocument() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear documento relacionado con fecha futura (inválido)
	relatedDoc := &models.RelatedDocument{}

	// Establecer valores
	b.setError(relatedDoc.SetDocumentType(constants.NotaRemisionElectronica))
	b.setError(relatedDoc.SetGenerationType(constants.ElectronicDocument))
	b.setError(relatedDoc.SetDocumentNumber("DA1E261A-BAD7-460F-AD15-04F2E281FC6A"))
	b.setError(relatedDoc.SetEmissionDate(utils.TimeNow().Add(24 * time.Hour))) // Mañana (inválido)

	// Añadir al documento
	if b.err == nil {
		relatedDocs := make([]interfaces.RelatedDocument, 0, 1)
		relatedDocs = append(relatedDocs, relatedDoc)
		b.setError(b.document.SetRelatedDocuments(relatedDocs))
	}

	return b
}

// AddOtherDocuments añade otros documentos predeterminados válidos
func (b *DTEBuilder) AddOtherDocuments() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Descripción y detalle para documento
	description := "Reference document"
	detail := "Reference document detail"

	// Crear documento
	otherDoc := &models.OtherDocument{}

	// Establecer valores
	b.setError(otherDoc.SetAssociatedDocument(constants.DocumentoEmisor))
	b.setError(otherDoc.SetDescription(description))
	b.setError(otherDoc.SetDetail(detail))

	// Añadir al documento
	if b.err == nil {
		otherDocs := make([]interfaces.OtherDocuments, 0, 1)
		otherDocs = append(otherDocs, otherDoc)
		b.setError(b.document.SetOtherDocuments(otherDocs))
	}

	return b
}

// AddMedicalDocument añade documento médico predeterminado válido
func (b *DTEBuilder) AddMedicalDocument() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear doctor
	doctor := &models.DoctorInfo{}

	// Establecer valores del doctor
	b.setError(doctor.SetName("Dr. John Smith"))
	b.setError(doctor.SetServiceType(1))

	nit := "12345678901234"
	b.setError(doctor.SetNIT(nit))

	// Crear documento médico
	medicalDoc := &models.OtherDocument{}

	// Establecer valores
	b.setError(medicalDoc.SetAssociatedDocument(constants.DocumentoMedico))
	b.setError(medicalDoc.SetDoctor(doctor))

	// Añadir o actualizar la lista de otros documentos
	var otherDocs []interfaces.OtherDocuments

	// Si ya existen otros documentos, agregarlos
	if existingDocs := b.document.GetOtherDocuments(); existingDocs != nil && len(existingDocs) > 0 {
		otherDocs = append(existingDocs, medicalDoc)
	} else {
		otherDocs = []interfaces.OtherDocuments{medicalDoc}
	}

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetOtherDocuments(otherDocs))
	}

	return b
}

// AddInvalidMedicalDocument añade documento médico inválido para testing
func (b *DTEBuilder) AddInvalidMedicalDocument() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear doctor
	doctor := &models.DoctorInfo{}

	// Establecer valores del doctor
	b.setError(doctor.SetName("Dr. John Smith"))
	b.setError(doctor.SetServiceType(1))

	nit := "12345678901234"
	b.setError(doctor.SetNIT(nit))

	// Descripción y detalle (inválidos para documento médico)
	description := "Invalid description for medical document"
	detail := "Invalid detail for medical document"

	// Crear documento médico
	medicalDoc := &models.OtherDocument{}

	// Establecer valores
	b.setError(medicalDoc.SetAssociatedDocument(constants.DocumentoMedico))
	b.setError(medicalDoc.SetDoctor(doctor))
	b.setError(medicalDoc.SetDescription(description))
	b.setError(medicalDoc.SetDetail(detail))

	// Añadir al documento
	if b.err == nil {
		otherDocs := make([]interfaces.OtherDocuments, 0, 1)
		otherDocs = append(otherDocs, medicalDoc)
		b.setError(b.document.SetOtherDocuments(otherDocs))
	}

	return b
}

// AddThirdPartySale añade venta de terceros predeterminada válida
func (b *DTEBuilder) AddThirdPartySale() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear venta de terceros
	thirdPartySale := &models.ThirdPartySale{}

	// Establecer valores
	b.setError(thirdPartySale.SetNIT("98765432101234"))
	b.setError(thirdPartySale.SetName("Third Party Company, Inc."))

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetThirdPartySale(thirdPartySale))
	}

	// Actualizar los items para incluir referencia al documento relacionado
	if b.err == nil {
		for i, itemInterface := range b.document.GetItems() {
			item, ok := itemInterface.(*models.Item)
			if ok {
				relatedDoc := fmt.Sprintf("THIRD-PARTY-DOC-%d", i+1)
				b.setError(item.SetRelatedDoc(&relatedDoc))
			}
		}
	}

	return b
}

// AddInvalidThirdPartySale añade venta de terceros inválida para testing
func (b *DTEBuilder) AddInvalidThirdPartySale() *DTEBuilder {
	// Si ya hay un error, no continuar
	if b.err != nil {
		return b
	}

	// Crear venta de terceros
	thirdPartySale := &models.ThirdPartySale{}

	// Establecer valores
	b.setError(thirdPartySale.SetNIT("98765432101234"))
	b.setError(thirdPartySale.SetName("Third Party Company, Inc."))

	// Asignar al documento
	if b.err == nil {
		b.setError(b.document.SetThirdPartySale(thirdPartySale))
	}

	// Actualizar SOLO EL PRIMER ITEM (mezcla de ventas propias y de terceros - inválido)
	if b.err == nil && len(b.document.GetItems()) > 0 {
		item, ok := b.document.GetItems()[0].(*models.Item)
		if ok {
			relatedDoc := "THIRD-PARTY-DOC-1"
			b.setError(item.SetRelatedDoc(&relatedDoc))
		}
	}

	return b
}

// BuildElectronicInvoice construye una factura electrónica válida
func (b *DTEBuilder) BuildElectronicInvoice() (*invoice_models.ElectronicInvoice, error) {
	// Construir un DTE base válido
	b.AddIdentification().
		AddIssuer().
		AddReceiver().
		AddItems().
		AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear una factura electrónica
	invoice := &invoice_models.ElectronicInvoice{
		DTEDocument:  baseDoc,
		InvoiceItems: make([]invoice_models.InvoiceItem, 0),
		InvoiceSummary: invoice_models.InvoiceSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a InvoiceItems
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)
		if ok {
			// Crear un InvoiceItem y asignar los valores
			invoiceItem := invoice_models.InvoiceItem{
				Item: baseItem,
			}

			// Calcular valores específicos de factura
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)
			ivaAmount := taxedAmount * 0.13

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			invoiceItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			invoiceItem.NonSubjectSale = *zeroAmountObj
			invoiceItem.ExemptSale = *zeroAmountObj
			invoiceItem.SuggestedPrice = *zeroAmountObj
			invoiceItem.NonTaxed = *zeroAmountObj

			ivaObj, err := financial.NewAmount(ivaAmount)
			if err != nil {
				return nil, err
			}
			invoiceItem.IVAItem = *ivaObj

			invoice.InvoiceItems = append(invoice.InvoiceItems, invoiceItem)
		}
	}

	// Agregar valores específicos del resumen de factura
	totalTaxed := invoice.GetSummary().GetTotalTaxed()
	totalIVA := totalTaxed * 0.13

	ivaAmountObj, err := financial.NewAmount(totalIVA)
	if err != nil {
		return nil, err
	}
	invoice.InvoiceSummary.TotalIva = *ivaAmountObj

	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	invoice.InvoiceSummary.TaxedDiscount = *zeroAmountObj
	invoice.InvoiceSummary.IVAPerception = *zeroAmountObj
	invoice.InvoiceSummary.IVARetention = *zeroAmountObj
	invoice.InvoiceSummary.IncomeRetention = *zeroAmountObj
	invoice.InvoiceSummary.BalanceInFavor = *zeroAmountObj

	// Validar la factura completa
	baseDTE := invoice.DTEDocument
	err = baseDTE.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := baseDTE.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return invoice, nil
}

// BuildInvalidElectronicInvoice construye una factura electrónica inválida
func (b *DTEBuilder) BuildInvalidElectronicInvoice() (*invoice_models.ElectronicInvoice, error) {
	// Construir un DTE base válido primero
	b.AddIdentification().
		AddIssuer().
		AddReceiver().
		AddItems().
		AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear una factura electrónica
	invoice := &invoice_models.ElectronicInvoice{
		DTEDocument:  baseDoc,
		InvoiceItems: make([]invoice_models.InvoiceItem, 0),
		InvoiceSummary: invoice_models.InvoiceSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a InvoiceItems con valores inválidos
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)
		if ok {
			// Crear un InvoiceItem y asignar los valores
			invoiceItem := invoice_models.InvoiceItem{
				Item: baseItem,
			}

			// Calcular valores específicos de factura
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)

			// IVA incorrecto (debería ser 13% de taxedAmount)
			ivaIncorrecto := taxedAmount * 0.20 // 20% en lugar de 13%

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			invoiceItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			invoiceItem.NonSubjectSale = *zeroAmountObj
			invoiceItem.ExemptSale = *zeroAmountObj
			invoiceItem.SuggestedPrice = *zeroAmountObj
			invoiceItem.NonTaxed = *zeroAmountObj

			ivaObj, err := financial.NewAmount(ivaIncorrecto)
			if err != nil {
				return nil, err
			}
			invoiceItem.IVAItem = *ivaObj

			invoice.InvoiceItems = append(invoice.InvoiceItems, invoiceItem)
		}
	}

	// Agregar valores específicos del resumen de factura con inconsistencias
	totalTaxed := invoice.GetSummary().GetTotalTaxed()

	// Total IVA incorrecto (inconsistente con los items)
	totalIVAIncorrecto := totalTaxed * 0.10 // 10% en lugar de 13%

	ivaAmountObj, err := financial.NewAmount(totalIVAIncorrecto)
	if err != nil {
		return nil, err
	}
	invoice.InvoiceSummary.TotalIva = *ivaAmountObj

	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	invoice.InvoiceSummary.TaxedDiscount = *zeroAmountObj
	invoice.InvoiceSummary.IVAPerception = *zeroAmountObj
	invoice.InvoiceSummary.IVARetention = *zeroAmountObj
	invoice.InvoiceSummary.IncomeRetention = *zeroAmountObj
	invoice.InvoiceSummary.BalanceInFavor = *zeroAmountObj

	return invoice, nil
}

// BuildCreditFiscalDocument construye un CCF válido
func (b *DTEBuilder) BuildCreditFiscalDocument() (*ccf_models.CreditFiscalDocument, error) {
	// Construir un DTE base válido con tipo CCF
	b.AddIdentification()

	// Modificar el tipo de documento a CCF
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.CCFElectronico))
	}

	// Continuar con la construcción
	b.AddIssuer().
		AddReceiverForCCF(). // Usar receptor específico para CCF
		AddItems().
		AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear un CCF
	ccf := &ccf_models.CreditFiscalDocument{
		DTEDocument: baseDoc,
		CreditItems: make([]ccf_models.CreditItem, 0),
		CreditSummary: ccf_models.CreditSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a CreditItems
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)
		if ok {
			// Crear un CreditItem y asignar los valores
			creditItem := ccf_models.CreditItem{
				Item: baseItem,
			}

			// Calcular valores específicos de CCF
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			creditItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			creditItem.NonSubjectSale = *zeroAmountObj
			creditItem.ExemptSale = *zeroAmountObj
			creditItem.SuggestedPrice = *zeroAmountObj
			creditItem.NonTaxed = *zeroAmountObj

			ccf.CreditItems = append(ccf.CreditItems, creditItem)
		}
	}

	// Agregar valores específicos del resumen de CCF
	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	ccf.CreditSummary.TaxedDiscount = *zeroAmountObj
	ccf.CreditSummary.IVAPerception = *zeroAmountObj
	ccf.CreditSummary.IVARetention = *zeroAmountObj
	ccf.CreditSummary.IncomeRetention = *zeroAmountObj
	ccf.CreditSummary.BalanceInFavor = *zeroAmountObj
	ccf.CreditSummary.ElectronicPaymentNumber = nil

	// Validar el CCF completo
	baseDTE := ccf.DTEDocument
	err = baseDTE.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := baseDTE.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return ccf, nil
}

// BuildInvalidCreditFiscalDocument construye un CCF inválido (receptor sin NRC)
func (b *DTEBuilder) BuildInvalidCreditFiscalDocument() (*ccf_models.CreditFiscalDocument, error) {
	// Construir un DTE base con tipo CCF pero receptor sin NRC (inválido para CCF)
	b.AddIdentification()

	// Modificar el tipo de documento a CCF
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.CCFElectronico))
	}

	// Continuar con la construcción pero usando un receptor sin NRC (inválido para CCF)
	b.AddIssuer().
		AddReceiverWithNoNRC(). // Receptor sin NRC, inválido para CCF
		AddItems().
		AddSummary()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear un CCF
	ccf := &ccf_models.CreditFiscalDocument{
		DTEDocument: baseDoc,
		CreditItems: make([]ccf_models.CreditItem, 0),
		CreditSummary: ccf_models.CreditSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a CreditItems
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)
		if ok {
			// Crear un CreditItem y asignar los valores
			creditItem := ccf_models.CreditItem{
				Item: baseItem,
			}

			// Calcular valores específicos de CCF
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			creditItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			creditItem.NonSubjectSale = *zeroAmountObj
			creditItem.ExemptSale = *zeroAmountObj
			creditItem.SuggestedPrice = *zeroAmountObj
			creditItem.NonTaxed = *zeroAmountObj

			ccf.CreditItems = append(ccf.CreditItems, creditItem)
		}
	}

	// Agregar valores específicos del resumen de CCF
	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	ccf.CreditSummary.TaxedDiscount = *zeroAmountObj
	ccf.CreditSummary.IVAPerception = *zeroAmountObj
	ccf.CreditSummary.IVARetention = *zeroAmountObj
	ccf.CreditSummary.IncomeRetention = *zeroAmountObj
	ccf.CreditSummary.BalanceInFavor = *zeroAmountObj
	ccf.CreditSummary.ElectronicPaymentNumber = nil

	return ccf, nil
}

// BuildCreditNote construye una nota de crédito válida
func (b *DTEBuilder) BuildCreditNote() (*credit_note_models.CreditNoteModel, error) {
	// Construir un DTE base válido con tipo Nota de Crédito
	b.AddIdentification()

	// Modificar el tipo de documento a Nota de Crédito
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.NotaCreditoElectronica))
	}

	// Continuar con la construcción
	b.AddIssuer().
		AddReceiver().
		AddItems().
		AddSummary().
		AddRelatedDocuments() // Las notas de crédito requieren documentos relacionados

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear una nota de crédito
	creditNote := &credit_note_models.CreditNoteModel{
		DTEDocument: baseDoc,
		CreditItems: make([]credit_note_models.CreditNoteItem, 0),
		CreditSummary: credit_note_models.CreditNoteSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a CreditNoteItems
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)

		if ok {
			// Crear un CreditNoteItem y asignar los valores
			creditNoteItem := credit_note_models.CreditNoteItem{
				Item: baseItem,
			}

			// Calcular valores específicos de Nota de Crédito
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			creditNoteItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			creditNoteItem.NonSubjectSale = *zeroAmountObj
			creditNoteItem.ExemptSale = *zeroAmountObj
			creditNoteItem.SuggestedPrice = *zeroAmountObj
			creditNoteItem.NonTaxed = *zeroAmountObj

			creditNote.CreditItems = append(creditNote.CreditItems, creditNoteItem)
		}
	}

	// Agregar valores específicos del resumen de Nota de Crédito
	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	creditNote.CreditSummary.TaxedDiscount = *zeroAmountObj
	creditNote.CreditSummary.IVAPerception = *zeroAmountObj
	creditNote.CreditSummary.IVARetention = *zeroAmountObj
	creditNote.CreditSummary.IncomeRetention = *zeroAmountObj
	creditNote.CreditSummary.BalanceInFavor = *zeroAmountObj

	// Validar la nota de crédito completa
	baseDTE := creditNote.DTEDocument
	err = baseDTE.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := baseDTE.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return creditNote, nil
}

// BuildInvalidCreditNote construye una nota de crédito inválida (sin documentos relacionados)
func (b *DTEBuilder) BuildInvalidCreditNote() (*credit_note_models.CreditNoteModel, error) {
	// Construir un DTE base con tipo Nota de Crédito pero sin documentos relacionados (inválido)
	b.AddIdentification()

	// Modificar el tipo de documento a Nota de Crédito
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.NotaCreditoElectronica))
	}

	// Continuar con la construcción pero sin añadir documentos relacionados (inválido)
	b.AddIssuer().
		AddReceiver().
		AddItems().
		AddSummary()
	// No añadimos documentos relacionados, lo que hace inválida la nota de crédito

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	// Obtener el documento base sin validaciones
	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	// Crear una nota de crédito
	creditNote := &credit_note_models.CreditNoteModel{
		DTEDocument: baseDoc,
		CreditItems: make([]credit_note_models.CreditNoteItem, 0),
		CreditSummary: credit_note_models.CreditNoteSummary{
			Summary: baseDoc.GetSummary().(*models.Summary),
		},
	}

	// Convertir los items genéricos a CreditNoteItems
	for _, item := range baseDoc.GetItems() {
		baseItem, ok := item.(*models.Item)
		if ok {
			// Crear un CreditNoteItem y asignar los valores
			creditNoteItem := credit_note_models.CreditNoteItem{
				Item: baseItem,
			}

			// Calcular valores específicos de Nota de Crédito
			taxedAmount := baseItem.GetQuantity() * baseItem.GetUnitPrice() * (1 - baseItem.GetDiscount()/100)

			// Establecer valores usando setters para atrapar errores
			amountObj, err := financial.NewAmount(taxedAmount)
			if err != nil {
				return nil, err
			}
			creditNoteItem.TaxedSale = *amountObj

			zeroAmountObj, err := financial.NewAmount(0)
			if err != nil {
				return nil, err
			}
			creditNoteItem.NonSubjectSale = *zeroAmountObj
			creditNoteItem.ExemptSale = *zeroAmountObj
			creditNoteItem.SuggestedPrice = *zeroAmountObj
			creditNoteItem.NonTaxed = *zeroAmountObj

			creditNote.CreditItems = append(creditNote.CreditItems, creditNoteItem)
		}
	}

	// Agregar valores específicos del resumen de Nota de Crédito
	zeroAmountObj, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	creditNote.CreditSummary.TaxedDiscount = *zeroAmountObj
	creditNote.CreditSummary.IVAPerception = *zeroAmountObj
	creditNote.CreditSummary.IVARetention = *zeroAmountObj
	creditNote.CreditSummary.IncomeRetention = *zeroAmountObj
	creditNote.CreditSummary.BalanceInFavor = *zeroAmountObj

	return creditNote, nil
}

// BuildRetentionDocumentWithPhysicalItems construye un documento de retención válido con items físicos
func (b *DTEBuilder) BuildRetentionDocumentWithPhysicalItems() (*retention_models.RetentionModel, error) {
	b.AddIdentification()

	// Modificar el tipo de documento a Retención
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.ComprobanteRetencionElectronico))
	}

	b.AddIssuer().
		AddReceiver()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	retentionDoc := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   make([]retention_models.RetentionItem, 0),
		RetentionSummary: &retention_models.RetentionSummary{},
	}

	// Crear items físicos de retención
	err = b.addPhysicalRetentionItems(retentionDoc)
	if err != nil {
		return nil, err
	}

	// Crear y asignar el resumen (requerido para items físicos)
	err = b.createRetentionSummary(retentionDoc)
	if err != nil {
		return nil, err
	}

	// Validar el documento de retención
	baseDTE := retentionDoc.DTEDocument
	err = baseDTE.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := baseDTE.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return retentionDoc, nil
}

// BuildRetentionDocumentWithElectronicItems construye un documento de retención con items electrónicos
func (b *DTEBuilder) BuildRetentionDocumentWithElectronicItems() (*retention_models.RetentionModel, error) {
	b.AddIdentification()

	// Modificar el tipo de documento a Retención
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.ComprobanteRetencionElectronico))
	}

	b.AddIssuer().
		AddReceiver()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	retentionDoc := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   make([]retention_models.RetentionItem, 0),
		RetentionSummary: &retention_models.RetentionSummary{},
	}

	err = b.addElectronicRetentionItems(retentionDoc)
	if err != nil {
		return nil, err
	}

	zeroAmount, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	retentionDoc.RetentionSummary.TotalIVARetention = *zeroAmount
	retentionDoc.RetentionSummary.TotalSubjectRetention = *zeroAmount

	// Validar el documento de retención
	baseDTE := retentionDoc.DTEDocument
	err = baseDTE.Validate()
	if err != nil {
		return nil, err
	}

	// Validar reglas de negocio
	dteErr := baseDTE.ValidateDTERules()
	if dteErr != nil {
		return nil, dteErr
	}

	return retentionDoc, nil
}

// BuildRetentionDocumentWithMixedItems construye un documento de retención con items mixtos (físicos y electrónicos)
func (b *DTEBuilder) BuildRetentionDocumentWithMixedItems() (*retention_models.RetentionModel, error) {
	b.AddIdentification()

	// Modificar el tipo de documento a Retención
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.ComprobanteRetencionElectronico))
	}

	b.AddIssuer().
		AddReceiver()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	retentionDoc := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   make([]retention_models.RetentionItem, 0),
		RetentionSummary: &retention_models.RetentionSummary{},
	}

	// Agregar una mezcla de items físicos y electrónicos
	err = b.addPhysicalRetentionItems(retentionDoc)
	if err != nil {
		return nil, err
	}
	err = b.addElectronicRetentionItems(retentionDoc)
	if err != nil {
		return nil, err
	}

	// Si hay items físicos, se necesita un resumen
	err = b.createRetentionSummary(retentionDoc)
	if err != nil {
		return nil, err
	}

	return retentionDoc, nil
}

// BuildInvalidRetentionDocument construye un documento de retención inválido con items físicos sin resumen
func (b *DTEBuilder) BuildInvalidRetentionDocument() (*retention_models.RetentionModel, error) {
	b.AddIdentification()

	// Modificar el tipo de documento a Retención
	identification, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identification != nil {
		b.setError(identification.SetDTEType(constants.ComprobanteRetencionElectronico))
	}

	b.AddIssuer().
		AddReceiver()

	// Si ya hay un error, no continuar
	if b.err != nil {
		return nil, b.err
	}

	baseDoc, err := b.BuildWithoutValidation()
	if err != nil {
		return nil, err
	}

	retentionDoc := &retention_models.RetentionModel{
		DTEDocument:      baseDoc,
		RetentionItems:   make([]retention_models.RetentionItem, 0),
		RetentionSummary: &retention_models.RetentionSummary{},
	}

	err = b.addPhysicalRetentionItems(retentionDoc)
	if err != nil {
		return nil, err
	}

	// ERROR: no creamos el resumen, lo que hace inválido el documento
	// para documentos con items físicos
	zeroAmount, err := financial.NewAmount(0)
	if err != nil {
		return nil, err
	}
	retentionDoc.RetentionSummary.TotalIVARetention = *zeroAmount
	retentionDoc.RetentionSummary.TotalSubjectRetention = *zeroAmount

	return retentionDoc, nil
}

// addPhysicalRetentionItems añade items físicos al documento de retención
func (b *DTEBuilder) addPhysicalRetentionItems(retentionDoc *retention_models.RetentionModel) error {
	// Crear un par de items físicos
	for i := 1; i <= 2; i++ {
		// Los items físicos requieren taxedAmount, ivaAmount, emissionDate y DTEType
		taxedAmount, err := financial.NewAmount(115.25 * float64(i))
		if err != nil {
			return err
		}

		ivaAmount, err := financial.NewAmount(15.00 * float64(i))
		if err != nil {
			return err
		}

		emissionDate, err := temporal.NewEmissionDate(utils.TimeNow().Add(-time.Hour * 24 * 30 * time.Duration(i))) // 1 o 2 meses atrás
		if err != nil {
			return err
		}

		dteType, err := document.NewDTEType(constants.CCFElectronico)
		if err != nil {
			return err
		}

		docType, err := document.NewOperationType(constants.PhysicalDocument)
		if err != nil {
			return err
		}

		docNumber, err := document.NewDocumentNumber(fmt.Sprintf("S221001%d", 340+i), constants.PhysicalDocument)
		if err != nil {
			return err
		}

		retentionCode, err := document.NewRetentionCode(constants.RetentionOnePercent)
		if err != nil {
			return err
		}

		retentionItem := retention_models.RetentionItem{
			Number:          *item.NewValidatedItemNumber(i),
			DocumentType:    *docType,
			DocumentNumber:  docNumber,
			Description:     fmt.Sprintf("Compra de suministros de oficina %d", i),
			RetentionAmount: *taxedAmount,
			RetentionIVA:    *ivaAmount,
			EmissionDate:    *emissionDate,
			DTEType:         *dteType,
			ReceptionCodeMH: *retentionCode,
		}

		retentionDoc.RetentionItems = append(retentionDoc.RetentionItems, retentionItem)
	}

	return nil
}

// addElectronicRetentionItems añade items electrónicos al documento de retención
func (b *DTEBuilder) addElectronicRetentionItems(retentionDoc *retention_models.RetentionModel) error {
	// Los números de items deben seguir la secuencia si ya hay items físicos
	startIdx := len(retentionDoc.RetentionItems) + 1

	for i := 0; i < 2; i++ {
		idx := startIdx + i

		docType, err := document.NewOperationType(constants.ElectronicDocument)
		if err != nil {
			return err
		}

		// Los items electrónicos necesitan un UUID válido como número de documento
		docNumber, err := document.NewDocumentNumber(fmt.Sprintf("FF54E9DB-79C3-42CE-B432-EC522C97EFB%d", i), constants.ElectronicDocument)
		if err != nil {
			return err
		}

		dteType, err := document.NewDTEType(constants.CCFElectronico)
		if err != nil {
			return err
		}

		retentionCode, err := document.NewRetentionCode(constants.RetentionThirteenPercent)
		if err != nil {
			return err
		}

		retentionItem := retention_models.RetentionItem{
			Number:          *item.NewValidatedItemNumber(idx),
			DTEType:         *dteType,
			DocumentType:    *docType,
			DocumentNumber:  docNumber,
			Description:     fmt.Sprintf("Servicio de consultoria electrónica %d", i+1),
			ReceptionCodeMH: *retentionCode,
		}

		retentionDoc.RetentionItems = append(retentionDoc.RetentionItems, retentionItem)
	}

	return nil
}

// createRetentionSummary crea y asigna un resumen para el documento de retención
func (b *DTEBuilder) createRetentionSummary(retentionDoc *retention_models.RetentionModel) error {
	// Calcular totales en base a los items físicos
	var totalSubjectRetention float64
	var totalIVARetention float64

	for _, item := range retentionDoc.RetentionItems {
		// Solo sumar los items físicos
		if item.DocumentType.GetValue() == constants.PhysicalDocument {
			totalSubjectRetention += item.RetentionAmount.GetValue()
			totalIVARetention += item.RetentionIVA.GetValue()
		}
	}

	totalSubject, err := financial.NewAmount(totalSubjectRetention)
	if err != nil {
		return err
	}

	totalIVA, err := financial.NewAmount(totalIVARetention)
	if err != nil {
		return err
	}

	retentionDoc.RetentionSummary.TotalSubjectRetention = *totalSubject
	retentionDoc.RetentionSummary.TotalIVARetention = *totalIVA

	return nil
}

// BuildInvalidationDocumentWithReplacement construye un documento de invalidación tipo 1 (con reemplazo)
func (b *DTEBuilder) BuildInvalidationDocumentWithReplacement() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.AddIdentification()

	identificationModel, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identificationModel != nil {
		b.setError(identificationModel.SetDTEType(constants.FacturaElectronica))
	}

	b.AddIssuer()

	invalidationDoc := &invalidation_models.InvalidationDocument{}

	invalidationDoc.Identification = identificationModel
	invalidationDoc.Issuer = b.document.GetIssuer().(*models.Issuer)

	invalidatedDoc, err := b.createInvalidatedDocument()
	if err != nil {
		return nil, err
	}

	replacementCode, err := identification.NewGenerationCode()
	if err != nil {
		return nil, err
	}
	invalidatedDoc.ReplacementCode = replacementCode

	invalidationDoc.Document = invalidatedDoc

	reason, err := b.createInvalidationReason(1)
	if err != nil {
		return nil, err
	}
	invalidationDoc.Reason = reason

	return invalidationDoc, nil
}

// BuildInvalidationDocumentWithAnnulment construye un documento de invalidación tipo 2 (anulación)
func (b *DTEBuilder) BuildInvalidationDocumentWithAnnulment() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.AddIdentification()

	identificationModel, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identificationModel != nil {
		b.setError(identificationModel.SetDTEType(constants.FacturaElectronica))
	}

	b.AddIssuer()

	invalidationDoc := &invalidation_models.InvalidationDocument{}

	invalidationDoc.Identification = identificationModel
	invalidationDoc.Issuer = b.document.GetIssuer().(*models.Issuer)

	invalidatedDoc, err := b.createInvalidatedDocument()
	if err != nil {
		return nil, err
	}

	invalidatedDoc.ReplacementCode = nil

	invalidationDoc.Document = invalidatedDoc

	reason, err := b.createInvalidationReason(2)
	if err != nil {
		return nil, err
	}
	invalidationDoc.Reason = reason

	return invalidationDoc, nil
}

// BuildInvalidationDocumentWithDefinitive construye un documento de invalidación tipo 3 (definitiva)
func (b *DTEBuilder) BuildInvalidationDocumentWithDefinitive() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.AddIdentification()

	identificationModel, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identificationModel != nil {
		b.setError(identificationModel.SetDTEType(constants.FacturaElectronica))
	}

	b.AddIssuer()

	invalidationDoc := &invalidation_models.InvalidationDocument{}

	invalidationDoc.Identification = identificationModel
	invalidationDoc.Issuer = b.document.GetIssuer().(*models.Issuer)

	invalidatedDoc, err := b.createInvalidatedDocument()
	if err != nil {
		return nil, err
	}

	replacementCode, err := identification.NewGenerationCode()
	if err != nil {
		return nil, err
	}
	invalidatedDoc.ReplacementCode = replacementCode

	invalidationDoc.Document = invalidatedDoc

	reason, err := b.createInvalidationReason(3)
	if err != nil {
		return nil, err
	}

	invalidReason := "Documento con errores graves que impiden su utilización"
	validatedReason, err := document.NewInvalidationReason(invalidReason)
	if err != nil {
		return nil, err
	}
	reason.Reason = validatedReason

	invalidationDoc.Reason = reason

	return invalidationDoc, nil
}

// BuildInvalidInvalidationDocument construye un documento de invalidación inválido (tipo 2 con código de reemplazo)
func (b *DTEBuilder) BuildInvalidInvalidationDocument() (*invalidation_models.InvalidationDocument, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.AddIdentification()

	identificationModel, ok := b.document.GetIdentification().(*models.Identification)
	if ok && identificationModel != nil {
		b.setError(identificationModel.SetDTEType(constants.FacturaElectronica))
	}

	b.AddIssuer()

	invalidationDoc := &invalidation_models.InvalidationDocument{}

	invalidationDoc.Identification = identificationModel
	invalidationDoc.Issuer = b.document.GetIssuer().(*models.Issuer)

	invalidatedDoc, err := b.createInvalidatedDocument()
	if err != nil {
		return nil, err
	}

	replacementCode, err := identification.NewGenerationCode()
	if err != nil {
		return nil, err
	}
	invalidatedDoc.ReplacementCode = replacementCode

	invalidationDoc.Document = invalidatedDoc

	reason, err := b.createInvalidationReason(2)
	if err != nil {
		return nil, err
	}
	invalidationDoc.Reason = reason

	return invalidationDoc, nil
}

// createInvalidatedDocument crea un documento invalidado para usar en los builders de invalidación
func (b *DTEBuilder) createInvalidatedDocument() (*invalidation_models.InvalidatedDocument, error) {
	docType, err := document.NewDTEType(constants.FacturaElectronica)
	if err != nil {
		return nil, err
	}

	generationCode, err := identification.NewGenerationCode()
	if err != nil {
		return nil, err
	}

	controlNumber, err := identification.NewControlNumber("DTE-01-00000000-000000000000001")
	if err != nil {
		return nil, err
	}

	receptionStamp := "2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM"

	emissionDate, err := temporal.NewEmissionDate(time.Now().Add(-24 * time.Hour))
	if err != nil {
		return nil, err
	}

	ivaAmount, err := financial.NewAmount(13.00)
	if err != nil {
		return nil, err
	}

	documentType, err := document.NewDTEType(constants.DUI)
	if err != nil {
		return nil, err
	}

	documentNumber, err := identification.NewDocumentNumber("01234567-8", constants.DUI)
	if err != nil {
		return nil, err
	}

	email, err := base.NewEmail("cliente@example.com")
	if err != nil {
		return nil, err
	}

	phone, err := base.NewPhone("22123456")
	if err != nil {
		return nil, err
	}

	name := "Cliente Ejemplo S.A. de C.V."

	return &invalidation_models.InvalidatedDocument{
		Type:           *docType,
		GenerationCode: *generationCode,
		ControlNumber:  *controlNumber,
		ReceptionStamp: receptionStamp,
		EmissionDate:   *emissionDate,
		IVAAmount:      ivaAmount,
		DocumentType:   documentType,
		DocumentNumber: documentNumber,
		Email:          email,
		Phone:          phone,
		Name:           &name,
	}, nil
}

// createInvalidationReason crea un motivo de invalidación para usar en los builders de invalidación
func (b *DTEBuilder) createInvalidationReason(invalidationType int) (*invalidation_models.InvalidationReason, error) {
	invalidationTypeObj, err := document.NewInvalidationType(invalidationType)
	if err != nil {
		return nil, err
	}

	responsibleDocType, err := document.NewDTETypeForReceiver(constants.DUI)
	if err != nil {
		return nil, err
	}

	responsibleDocNum, err := identification.NewDocumentNumber("01234567-8", constants.DUI)
	if err != nil {
		return nil, err
	}

	requestorDocType, err := document.NewDTETypeForReceiver(constants.DUI)
	if err != nil {
		return nil, err
	}

	requestorDocNum, err := identification.NewDocumentNumber("98765432-1", constants.DUI)
	if err != nil {
		return nil, err
	}

	return &invalidation_models.InvalidationReason{
		Type:               *invalidationTypeObj,
		ResponsibleName:    "Juan Responsable",
		ResponsibleDocType: *responsibleDocType,
		ResponsibleDocNum:  *responsibleDocNum,
		RequesterName:      "Ana Solicitante",
		RequesterDocType:   *requestorDocType,
		RequesterDocNum:    *requestorDocNum,
	}, nil
}
