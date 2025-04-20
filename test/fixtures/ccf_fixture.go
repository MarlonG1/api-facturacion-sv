package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// CreateDefaultCreditItem crea un ítem de CCF predeterminado válido
func CreateDefaultCreditItem(index int) structs.CreditItemRequest {
	code := "CCF" + string(rune(65+index))

	return structs.CreditItemRequest{
		ItemRequest: structs.ItemRequest{
			Number:      index + 1,
			Type:        1, // Producto
			Description: "Producto CCF " + string(rune(65+index)),
			Quantity:    15,
			UnitMeasure: 59, // Unidades
			UnitPrice:   10.0,
			Discount:    0,
			Code:        &code,
			Taxes:       []string{"20"}, // Código IVA
		},
		NonSubjectSale: 0,
		ExemptSale:     0,
		TaxedSale:      150.0, // Cantidad * Precio unitario
		SuggestedPrice: 0,
		NonTaxed:       0,
	}
}

// CreateDefaultCreditSummary crea un resumen de CCF predeterminado válido
func CreateDefaultCreditSummary() *structs.CreditSummaryRequest {
	return &structs.CreditSummaryRequest{
		SummaryRequest: structs.SummaryRequest{
			TotalNonSubject:    0,
			TotalExempt:        0,
			TotalTaxed:         300.0,
			SubTotal:           300.0,
			NonSubjectDiscount: 0,
			ExemptDiscount:     0,
			DiscountPercentage: 0,
			TotalDiscount:      0,
			SubTotalSales:      300.0,
			TotalOperation:     300.0,
			TotalNonTaxed:      0,
			TotalToPay:         300.0,
			OperationCondition: 1, // Contado
			Taxes: []structs.TaxRequest{
				{
					Code:        "20", // Código IVA
					Description: "IVA",
					Value:       39.0, // 13% del monto gravado
				},
			},
			PaymentTypes: []structs.PaymentRequest{
				{
					Code:   "01", // Efectivo
					Amount: 300.0,
				},
			},
		},
		TaxedDiscount:   0,
		IVAPerception:   0,
		IVARetention:    0,
		IncomeRetention: 0,
		BalanceInFavor:  0,
	}
}

// CreateDefaultCreditFiscalRequest crea una solicitud de CCF predeterminada válida
func CreateDefaultCreditFiscalRequest() *structs.CreateCreditFiscalRequest {
	items := []structs.CreditItemRequest{
		CreateDefaultCreditItem(1),
		CreateDefaultCreditItem(2),
	}

	return &structs.CreateCreditFiscalRequest{
		Items:     items,
		Receiver:  CreateDefaultReceiverWithoutDocsFields(),
		ModelType: 1, // Modelo normal
		Summary:   CreateDefaultCreditSummary(),
	}
}

// CreateCreditFiscalWithInvalidItems crea una solicitud de CCF con ítems inválidos
func CreateCreditFiscalWithInvalidItems() *structs.CreateCreditFiscalRequest {
	req := CreateDefaultCreditFiscalRequest()
	item := CreateDefaultCreditItem(0)
	item.Type = 99 // Tipo inválido
	req.Items = []structs.CreditItemRequest{item}
	return req
}

// CreateCreditFiscalWithNonSubjectSale crea una solicitud de CCF con venta no sujeta (inválido)
func CreateCreditFiscalWithNonSubjectSale() *structs.CreateCreditFiscalRequest {
	req := CreateDefaultCreditFiscalRequest()
	item := CreateDefaultCreditItem(0)
	item.NonSubjectSale = 50.0 // CCF no puede incluir ventas no sujetas
	req.Items = []structs.CreditItemRequest{item}
	return req
}

// CreateCCFRequestWithAllOptionalFields crea una solicitud de factura con todos los campos opcionales
func CreateCCFRequestWithAllOptionalFields() *structs.CreateCreditFiscalRequest {
	req := CreateDefaultCreditFiscalRequest()
	req.Extension = CreateDefaultExtension()
	req.ThirdPartySale = CreateDefaultThirdPartySale()
	req.RelatedDocs = []structs.RelatedDocRequest{CreateDefaultRelatedDocument()}
	req.OtherDocs = []structs.OtherDocRequest{CreateDefaultOtherDocument()}
	req.Appendixes = []structs.AppendixRequest{CreateDefaultAppendix()}
	return req
}
