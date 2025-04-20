package fixtures

import (
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// CreateAddressWithEmptyFields crea una dirección con campos vacíos
func CreateAddressWithEmptyFields() *structs.AddressRequest {
	return &structs.AddressRequest{
		Department:   "",
		Municipality: "",
		Complement:   "",
	}
}

// CreateAddressWithInvalidMunicipality crea una dirección con municipio inválido
func CreateAddressWithInvalidMunicipality() *structs.AddressRequest {
	address := CreateDefaultAddress()
	address.Municipality = "99" // Código de municipio inválido
	return address
}

// CreateReceiverWithInvalidEmail crea un receptor con email inválido
func CreateReceiverWithInvalidEmail() *structs.ReceiverRequest {
	receiver := CreateDefaultReceiver()
	invalidEmail := "not-an-email"
	receiver.Email = &invalidEmail
	return receiver
}

// CreateReceiverWithoutRequiredFields crea un receptor sin campos requeridos
func CreateReceiverWithoutRequiredFields() *structs.ReceiverRequest {
	receiver := CreateDefaultReceiver()
	receiver.Name = nil
	receiver.NRC = nil
	return receiver
}

// CreateExtensionWithMissingFields crea una extensión con campos faltantes
func CreateExtensionWithMissingFields() *structs.ExtensionRequest {
	ext := CreateDefaultExtension()
	ext.DeliveryName = ""
	ext.ReceiverName = ""
	return ext
}

// CreateAppendixWithInvalidField crea un apéndice con campo inválido
func CreateAppendixWithInvalidField() structs.AppendixRequest {
	appendix := CreateDefaultAppendix()
	appendix.Field = ""
	return appendix
}

// CreatePaymentWithInvalidCode crea un pago con código inválido
func CreatePaymentWithInvalidCode() structs.PaymentRequest {
	payment := CreateDefaultPayment()
	payment.Code = "100" // Código de pago inválido
	return payment
}

// CreateThirdPartySaleWithEmptyNIT crea una venta a terceros con NIT vacío
func CreateThirdPartySaleWithEmptyNIT() *structs.ThirdPartySaleRequest {
	sale := CreateDefaultThirdPartySale()
	sale.NIT = ""
	return sale
}
