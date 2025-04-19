package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// CreateDefaultCreditNoteItem crea un ítem de nota de crédito predeterminado válido
func CreateDefaultCreditNoteItem(index int) structs.CreditNoteItemRequest {
	code := "CN" + string(rune(65+index))

	return structs.CreditNoteItemRequest{
		ItemRequest: structs.ItemRequest{
			Number:      index + 1,
			Type:        1, // Producto
			Description: "Devolución Producto " + string(rune(65+index)),
			Quantity:    5,
			UnitMeasure: 59, // Unidades
			UnitPrice:   5.0,
			Discount:    0,
			Code:        &code,
			Taxes:       []string{"20"}, // Código IVA
		},
		NonSubjectSale: 0,
		ExemptSale:     0,
		TaxedSale:      25.0, // Cantidad * Precio unitario
		SuggestedPrice: 0,
		NonTaxed:       0,
	}
}

// CreateDefaultCreditNoteSummary crea un resumen de nota de crédito predeterminado válido
func CreateDefaultCreditNoteSummary() *structs.CreditNoteSummaryRequest {
	return &structs.CreditNoteSummaryRequest{
		SummaryRequest: structs.SummaryRequest{
			TotalNonSubject:    0,
			TotalExempt:        0,
			TotalTaxed:         50.0,
			SubTotal:           50.0,
			NonSubjectDiscount: 0,
			ExemptDiscount:     0,
			DiscountPercentage: 0,
			TotalDiscount:      0,
			SubTotalSales:      50.0,
			TotalOperation:     50.0,
			TotalNonTaxed:      0,
			TotalToPay:         1, // Por convención en Notas de Crédito
			OperationCondition: 1, // Contado
			Taxes: []structs.TaxRequest{
				{
					Code:        "20", // Código IVA
					Description: "IVA",
					Value:       6.5, // 13% del monto gravado
				},
			},
			PaymentTypes: []structs.PaymentRequest{},
		},
		TaxedDiscount:   0,
		IVAPerception:   0,
		IVARetention:    0,
		IncomeRetention: 0,
		BalanceInFavor:  0,
	}
}

// CreateDefaultCreditNoteRequest crea una solicitud de nota de crédito predeterminada válida
func CreateDefaultCreditNoteRequest() *structs.CreateCreditNoteRequest {
	items := []structs.CreditNoteItemRequest{
		CreateDefaultCreditNoteItem(1),
		CreateDefaultCreditNoteItem(2),
	}

	return &structs.CreateCreditNoteRequest{
		Items:     items,
		Receiver:  CreateDefaultReceiver(),
		ModelType: constants.ModeloFacturacionPrevio, // Modelo normal
		Summary:   CreateDefaultCreditNoteSummary(),
		RelatedDocs: []structs.RelatedDocRequest{
			CreateDefaultRelatedDocument(),
		},
	}
}

// CreateDefaultCreditNoteExtension crea una extensión predeterminada válida
func CreateDefaultCreditNoteExtension() *structs.ExtensionRequest {
	observation := "Observación de prueba"

	return &structs.ExtensionRequest{
		DeliveryName:     "Juan Pérez",
		DeliveryDocument: "12345678-9",
		ReceiverName:     "Ana López",
		ReceiverDocument: "98765432-1",
		Observation:      &observation,
	}
}

// CreateCreditNoteWithInvalidItems crea una solicitud de nota de crédito con ítems inválidos
func CreateCreditNoteWithInvalidItems() *structs.CreateCreditNoteRequest {
	req := CreateDefaultCreditNoteRequest()
	item := CreateDefaultCreditNoteItem(0)
	item.Type = 99 // Tipo inválido
	req.Items = []structs.CreditNoteItemRequest{item}
	return req
}

// CreateCreditNoteWithoutRelatedDocs crea una solicitud de nota de crédito sin documentos relacionados
func CreateCreditNoteWithoutRelatedDocs() *structs.CreateCreditNoteRequest {
	req := CreateDefaultCreditNoteRequest()
	req.RelatedDocs = nil
	return req
}

// CreateCreditNoteRequestWithAllOptionalFields crea una solicitud de factura con todos los campos opcionales
func CreateCreditNoteRequestWithAllOptionalFields() *structs.CreateCreditNoteRequest {
	req := CreateDefaultCreditNoteRequest()
	req.Extension = CreateDefaultCreditNoteExtension()
	req.ThirdPartySale = CreateDefaultThirdPartySale()
	req.RelatedDocs = []structs.RelatedDocRequest{CreateDefaultRelatedDocument()}
	req.OtherDocs = []structs.OtherDocRequest{CreateDefaultOtherDocument()}
	req.Appendixes = []structs.AppendixRequest{CreateDefaultAppendix()}
	return req
}
