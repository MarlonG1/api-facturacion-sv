package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// CreateDefaultInvoiceItem crea un ítem de factura predeterminado válido
func CreateDefaultInvoiceItem(index int) structs.InvoiceItemRequest {
	code := "COD" + string(rune(65+index))

	return structs.InvoiceItemRequest{
		ItemRequest: structs.ItemRequest{
			Number:      index + 1,
			Type:        1, // Producto
			Description: "Producto " + string(rune(65+index)),
			Quantity:    10,
			UnitMeasure: 59, // Unidades
			UnitPrice:   5.0,
			Discount:    0,
			Code:        &code,
			Taxes:       []string{"20"}, // Código IVA
		},
		NonSubjectSale: 0,
		ExemptSale:     0,
		TaxedSale:      50.0, // Cantidad * Precio unitario
		SuggestedPrice: 0,
		NonTaxed:       0,
		IVAItem:        6.5, // TaxedSale * 0.13
	}
}

// CreateInvoiceItemWithInvalidType crea un ítem de factura con tipo inválido
func CreateInvoiceItemWithInvalidType() structs.InvoiceItemRequest {
	item := CreateDefaultInvoiceItem(0)
	item.Type = 99 // Tipo inválido
	return item
}

// CreateInvoiceItemWithNegativeQuantity crea un ítem de factura con cantidad negativa
func CreateInvoiceItemWithNegativeQuantity() structs.InvoiceItemRequest {
	item := CreateDefaultInvoiceItem(0)
	item.Quantity = -1 // Cantidad negativa
	return item
}

// CreateInvoiceItemWithInvalidIVA crea un ítem de factura con IVA inválido
func CreateInvoiceItemWithInvalidIVA() structs.InvoiceItemRequest {
	item := CreateDefaultInvoiceItem(0)
	item.IVAItem = 100 // Monto IVA inválido
	return item
}

// CreateDefaultInvoiceSummary crea un resumen de factura predeterminado válido
func CreateDefaultInvoiceSummary() *structs.InvoiceSummaryRequest {
	return &structs.InvoiceSummaryRequest{
		SummaryRequest: structs.SummaryRequest{
			TotalNonSubject:    0,
			TotalExempt:        0,
			TotalTaxed:         100.0,
			SubTotal:           100.0,
			NonSubjectDiscount: 0,
			ExemptDiscount:     0,
			DiscountPercentage: 0,
			TotalDiscount:      0,
			SubTotalSales:      100.0,
			TotalOperation:     100.0,
			TotalNonTaxed:      0,
			TotalToPay:         100.0,
			OperationCondition: 1, // Contado
			Taxes: []structs.TaxRequest{
				{
					Code:        "20", // Código IVA
					Description: "IVA",
					Value:       13.0, // 13% del monto gravado
				},
			},
			PaymentTypes: []structs.PaymentRequest{
				{
					Code:   "01", // Efectivo
					Amount: 100.0,
				},
			},
		},
		TaxedDiscount:   0,
		IVAPerception:   0,
		IVARetention:    0,
		IncomeRetention: 0,
		TotalIVA:        13.0,
		BalanceInFavor:  0,
	}
}

// CreateInvoiceSummaryWithInvalidTaxCode crea un resumen de factura con código de impuesto inválido
func CreateInvoiceSummaryWithInvalidTaxCode() *structs.InvoiceSummaryRequest {
	summary := CreateDefaultInvoiceSummary()
	summary.Taxes[0].Code = "99" // Código de impuesto inválido
	return summary
}

// CreateInvoiceSummaryWithMismatchedTotals crea un resumen de factura con totales inconsistentes
func CreateInvoiceSummaryWithMismatchedTotals() *structs.InvoiceSummaryRequest {
	summary := CreateDefaultInvoiceSummary()
	summary.TotalToPay = 50.0 // No coincide con TotalOperation
	return summary
}

// CreateDefaultInvoiceRequest crea una solicitud de factura predeterminada válida
func CreateDefaultInvoiceRequest() *structs.CreateInvoiceRequest {
	items := []structs.InvoiceItemRequest{
		CreateDefaultInvoiceItem(0),
		CreateDefaultInvoiceItem(1),
	}

	return &structs.CreateInvoiceRequest{
		Items:     items,
		Receiver:  CreateDefaultReceiver(),
		ModelType: 1, // Modelo normal
		Summary:   CreateDefaultInvoiceSummary(),
	}
}

// CreateInvoiceRequestWithInvalidItems crea una solicitud de factura con ítems inválidos
func CreateInvoiceRequestWithInvalidItems() *structs.CreateInvoiceRequest {
	req := CreateDefaultInvoiceRequest()
	req.Items = []structs.InvoiceItemRequest{
		CreateInvoiceItemWithInvalidType(),
	}
	return req
}

// CreateInvoiceRequestWithInvalidSummary crea una solicitud de factura con resumen inválido
func CreateInvoiceRequestWithInvalidSummary() *structs.CreateInvoiceRequest {
	req := CreateDefaultInvoiceRequest()
	req.Summary = CreateInvoiceSummaryWithMismatchedTotals()
	return req
}

// CreateInvoiceRequestWithInvalidReceiver crea una solicitud de factura con receptor inválido
func CreateInvoiceRequestWithInvalidReceiver() *structs.CreateInvoiceRequest {
	req := CreateDefaultInvoiceRequest()
	req.Receiver = CreateReceiverWithInvalidEmail()
	return req
}

// CreateInvoiceRequestWithAllOptionalFields crea una solicitud de factura con todos los campos opcionales
func CreateInvoiceRequestWithAllOptionalFields() *structs.CreateInvoiceRequest {
	req := CreateDefaultInvoiceRequest()
	req.Extension = CreateDefaultExtension()
	req.ThirdPartySale = CreateDefaultThirdPartySale()
	req.RelatedDocs = []structs.RelatedDocRequest{CreateDefaultRelatedDocument()}
	req.OtherDocs = []structs.OtherDocRequest{CreateDefaultOtherDocument()}
	req.Appendixes = []structs.AppendixRequest{CreateDefaultAppendix()}
	return req
}
