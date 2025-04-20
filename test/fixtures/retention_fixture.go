package fixtures

import "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"

// CreatePhysicalDocumentsRetentionRequest crea una solicitud de retención por defecto con documentos físicos
func CreatePhysicalDocumentsRetentionRequest() *structs.CreateRetentionRequest {
	// Crear valores por defecto para montos y fechas
	taxedAmount1 := 115.25
	ivaAmount1 := 15.00
	emissionDate1 := "2025-03-20"
	dteType1 := "03" // CCF

	taxedAmount2 := 226.50
	ivaAmount2 := 29.43
	emissionDate2 := "2025-03-22"
	dteType2 := "03" // CCF

	// Crear ítems físicos
	items := []structs.RetentionItem{
		{
			DocumentType:   1, // Físico
			DocumentNumber: "S221001345",
			Description:    "Compra de suministros de oficina",
			RetentionCode:  "22", // IVA 1%
			TaxedAmount:    &taxedAmount1,
			IvaAmount:      &ivaAmount1,
			EmissionDate:   &emissionDate1,
			DTEType:        &dteType1,
		},
		{
			DocumentType:   1, // Físico
			DocumentNumber: "S221001346",
			Description:    "Servicio de limpieza",
			RetentionCode:  "C4", // IVA 13%
			TaxedAmount:    &taxedAmount2,
			IvaAmount:      &ivaAmount2,
			EmissionDate:   &emissionDate2,
			DTEType:        &dteType2,
		},
	}

	// Crear receptor
	docType := "36" // NIT
	docNumber := "06141804941035"
	nrc := "123456"
	name := "Empresa Servicios Generales, S.A. de C.V."
	commercialName := "ServiGeneral"
	activityCode := "46900"
	activityDesc := "Venta al por mayor de otros productos"
	phone := "22123456"
	email := "info@gmail.com"

	receiver := &structs.ReceiverRequest{
		DocumentType:   &docType,
		DocumentNumber: &docNumber,
		NRC:            &nrc,
		Name:           &name,
		CommercialName: &commercialName,
		ActivityCode:   &activityCode,
		ActivityDesc:   &activityDesc,
		Address: &structs.AddressRequest{
			Department:   "06",
			Municipality: "20",
			Complement:   "Colonia Escalón, Calle La Reforma #123, San Salvador",
		},
		Phone: &phone,
		Email: &email,
	}

	// Crear resumen
	summary := &structs.RetentionSummary{
		TotalRetentionAmount: 341.75,
		TotalRetentionIVA:    44.43,
	}

	return &structs.CreateRetentionRequest{
		Items:    items,
		Receiver: receiver,
		Summary:  summary,
	}
}

// CreateElectronicDocumentsRetentionRequest crea una solicitud de retención por defecto con documentos electrónicos
func CreateElectronicDocumentsRetentionRequest() *structs.CreateRetentionRequest {
	// Crear ítems electrónicos
	items := []structs.RetentionItem{
		{
			DocumentType:   2, // Electrónico
			DocumentNumber: "FF54E9DB-79C3-42CE-B432-EC522C97EFB9",
			Description:    "Compra de equipos informáticos",
			RetentionCode:  "22", // IVA 1%
		},
		{
			DocumentType:   2, // Electrónico
			DocumentNumber: "AD54E9BB-79A3-42AE-B432-EC522C97EFB7",
			Description:    "Mantenimiento de servidores",
			RetentionCode:  "C4", // IVA 13%
		},
	}

	// Crear receptor
	docType := "36" // NIT
	docNumber := "06141804941035"
	nrc := "123456"
	name := "Empresa Servicios Generales, S.A. de C.V."
	commercialName := "ServiGeneral"
	activityCode := "46900"
	activityDesc := "Venta al por mayor de otros productos"
	phone := "22123456"
	email := "info@gmail.com"

	receiver := &structs.ReceiverRequest{
		DocumentType:   &docType,
		DocumentNumber: &docNumber,
		NRC:            &nrc,
		Name:           &name,
		CommercialName: &commercialName,
		ActivityCode:   &activityCode,
		ActivityDesc:   &activityDesc,
		Address: &structs.AddressRequest{
			Department:   "06",
			Municipality: "20",
			Complement:   "Colonia Escalón, Calle La Reforma #123, San Salvador",
		},
		Phone: &phone,
		Email: &email,
	}

	// Crear extension opcional
	observation := "Retención por servicios tecnológicos primer trimestre"
	extension := &structs.ExtensionRequest{
		Observation:      &observation,
		DeliveryName:     "Juan Carlos Martínez",
		DeliveryDocument: "04567890-1",
		ReceiverName:     "Ana María López",
		ReceiverDocument: "12345678-9",
	}

	return &structs.CreateRetentionRequest{
		Items:     items,
		Receiver:  receiver,
		Extension: extension,
	}
}

// CreateMixedDocumentsRetentionRequest crea una solicitud de retención por defecto con documentos mixtos
func CreateMixedDocumentsRetentionRequest() *structs.CreateRetentionRequest {
	// Crear valores por defecto para montos y fechas para el ítem físico
	taxedAmount := 450.00
	ivaAmount := 58.50
	emissionDate := "2025-03-15"
	dteType := "03" // CCF

	// Crear ítems mixtos (físico y electrónico)
	items := []structs.RetentionItem{
		{
			DocumentType:   1, // Físico
			DocumentNumber: "S221001347",
			Description:    "Consultoría financiera",
			RetentionCode:  "C9", // Otros casos
			TaxedAmount:    &taxedAmount,
			IvaAmount:      &ivaAmount,
			EmissionDate:   &emissionDate,
			DTEType:        &dteType,
		},
		{
			DocumentType:   2, // Electrónico
			DocumentNumber: "FF32E9DB-79C3-42CE-B432-EC522C97EFB2",
			Description:    "Servicios de auditoría",
			RetentionCode:  "C4", // IVA 13%
		},
	}

	// Crear receptor
	docType := "36" // NIT
	docNumber := "06141804941035"
	nrc := "123456"
	name := "Empresa Servicios Generales, S.A. de C.V."
	commercialName := "ServiGeneral"
	activityCode := "46900"
	activityDesc := "Venta al por mayor de otros productos"
	phone := "22123456"
	email := "info@gmail.com"

	receiver := &structs.ReceiverRequest{
		DocumentType:   &docType,
		DocumentNumber: &docNumber,
		NRC:            &nrc,
		Name:           &name,
		CommercialName: &commercialName,
		ActivityCode:   &activityCode,
		ActivityDesc:   &activityDesc,
		Address: &structs.AddressRequest{
			Department:   "06",
			Municipality: "20",
			Complement:   "Colonia Escalón, Calle La Reforma #123, San Salvador",
		},
		Phone: &phone,
		Email: &email,
	}

	// Crear apéndices opcionales
	appendixes := []structs.AppendixRequest{
		{
			Field: "nota_interna",
			Label: "Nota interna",
			Value: "Retención realizada según contrato marco",
		},
	}

	return &structs.CreateRetentionRequest{
		Items:      items,
		Receiver:   receiver,
		Appendixes: appendixes,
	}
}

// CreateRetentionRequestWithAllOptionalFields crea una solicitud de retención con todos los campos opcionales
func CreateRetentionRequestWithAllOptionalFields() *structs.CreateRetentionRequest {
	req := CreatePhysicalDocumentsRetentionRequest()

	// Agregar extensión
	observation := "Observación de prueba para retención"
	req.Extension = &structs.ExtensionRequest{
		Observation:      &observation,
		DeliveryName:     "Juan Pérez",
		DeliveryDocument: "12345678-9",
		ReceiverName:     "María González",
		ReceiverDocument: "98765432-1",
	}

	// Agregar apéndices
	req.Appendixes = []structs.AppendixRequest{
		{
			Field: "campo_adicional",
			Label: "Información adicional",
			Value: "Esta es información adicional para la retención",
		},
		{
			Field: "referencia_interna",
			Label: "Referencia interna",
			Value: "REF-2023-001",
		},
	}

	return req
}
