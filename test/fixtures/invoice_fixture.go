package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	respStructs "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
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

// CreateExpectedInvoiceResponse crea una respuesta esperada de factura para pruebas
func CreateExpectedInvoiceResponse() *respStructs.InvoiceDTEResponse {
	// Campos de identificación
	identificacion := &respStructs.DTEIdentification{
		Version:          1,
		Ambiente:         "00",
		TipoDte:          "01",
		NumeroControl:    "DTE-01-00000000-000000000000001",
		CodigoGeneracion: "FF54E9DB-79C3-42CE-B432-EC522C97EFB9", // Será ignorado en las comparaciones
		TipoModelo:       1,
		TipoOperacion:    1,
		TipoContingencia: nil,
		MotivoContin:     nil,
		FecEmi:           "2025-04-18", // Sera ignroado en las comparaciones - fecha actual en formato YYYY-MM-DD
		HorEmi:           "15:30:00",   // Sera ignorado en las comparaciones - hora actual en formato HH:MM:SS
		TipoMoneda:       "USD",
	}

	// Datos del emisor
	emisor := respStructs.DTEIssuer{
		NIT:                 "11111111111111",
		NRC:                 "1111111",
		Nombre:              "EMPRESA DE PRUEBAS SA DE CV",
		CodActividad:        "11111",
		DescActividad:       "Venta al por mayor de otros productos",
		NombreComercial:     utils.ToStringPointer("EJEMPLO SA"),
		TipoEstablecimiento: "02",
		Direccion: respStructs.DTEAddress{
			Departamento: "06",
			Municipio:    "20",
			Complemento:  "BOULEVARD SANTA ELENA SUR, SANTA TECLA",
		},
		Telefono:        "22567890",
		Correo:          "email@gmail.com",
		CodEstableMH:    nil,
		CodEstable:      utils.ToStringPointer("C001"),
		CodPuntoVentaMH: nil,
		CodPuntoVenta:   nil,
	}

	// Datos del receptor
	receptor := respStructs.InvoiceReceiver{
		Nombre:        utils.ToStringPointer("Empresa Servicios Generales, S.A. de C.V."),
		TipoDocumento: utils.ToStringPointer("36"),
		NumDocumento:  utils.ToStringPointer("06141804941035"),
		NRC:           utils.ToStringPointer("123456"),
		CodActividad:  utils.ToStringPointer("46900"),
		DescActividad: utils.ToStringPointer("Venta al por mayor de otros productos"),
		Direccion: &respStructs.DTEAddress{
			Departamento: "06",
			Municipio:    "20",
			Complemento:  "Colonia Escalón, Calle La Reforma #123, San Salvador",
		},
		Telefono: utils.ToStringPointer("22123456"),
		Correo:   utils.ToStringPointer("empresa@example.com"),
	}

	// Ítems del documento
	items := []respStructs.InvoiceItem{
		{
			NumItem:      1,
			TipoItem:     1,
			Codigo:       utils.ToStringPointer("CODA"),
			Descripcion:  "Producto A",
			Cantidad:     10,
			UniMedida:    59,
			PrecioUni:    5,
			MontoDescu:   0,
			VentaNoSuj:   0,
			VentaExenta:  0,
			VentaGravada: 50,
			Tributos:     []string{"20"},
			PSV:          0,
			NoGravado:    0,
			IvaItem:      6.5,
		},
		{
			NumItem:      2,
			TipoItem:     1,
			Codigo:       utils.ToStringPointer("CODB"),
			Descripcion:  "Producto B",
			Cantidad:     10,
			UniMedida:    59,
			PrecioUni:    5,
			MontoDescu:   0,
			VentaNoSuj:   0,
			VentaExenta:  0,
			VentaGravada: 50,
			Tributos:     []string{"20"},
			PSV:          0,
			NoGravado:    0,
			IvaItem:      6.5,
		},
	}

	// Resumen
	resumen := &respStructs.InvoiceSummary{
		TotalNoSuj:          0,
		TotalExenta:         0,
		TotalGravada:        100,
		SubTotalVentas:      100,
		DescuNoSuj:          0,
		DescuExenta:         0,
		DescuGravada:        0,
		PorcentajeDescuento: 0,
		TotalDescu:          0,
		SubTotal:            100,
		ReteRenta:           0,
		IvaRete1:            0,
		IvaPerci1:           nil,
		MontoTotalOperacion: 100,
		TotalNoGravado:      0,
		TotalPagar:          100,
		TotalLetras:         "CIEN DÓLARES",
		TotalIva:            13.0,
		SaldoFavor:          0,
		CondicionOperacion:  1,
		Pagos: []respStructs.DTEPayment{
			{
				Codigo:     "01",
				MontoPago:  100,
				Referencia: utils.ToStringPointer(""),
				Plazo:      nil,
				Periodo:    nil,
			},
		},
		NumPagoElectronico: nil,
	}

	// Construir respuesta completa
	return &respStructs.InvoiceDTEResponse{
		Identificacion:  identificacion,
		Emisor:          emisor,
		Receptor:        receptor,
		CuerpoDocumento: items,
		Resumen:         resumen,
		// Otros campos quedan como nil al no estar en el fixture básico
	}
}
