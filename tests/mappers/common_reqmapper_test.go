package mappers

import (
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/tests"
	"github.com/MarlonG1/api-facturacion-sv/tests/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestCommonMappers(t *testing.T) {
	test.TestMain(t)

	// Test para MapCommonRequestAddress
	t.Run("TestMapCommonRequestAddress", func(t *testing.T) {
		// Caso válido
		address := fixtures.CreateDefaultAddress()
		result, err := common.MapCommonRequestAddress(*address)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, address.Department, result.Department.GetValue())
		assert.Equal(t, address.Municipality, result.Municipality.GetValue())
		assert.Equal(t, address.Complement, result.Complement.GetValue())

		// Caso inválido: campos vacíos
		addressInvalid := fixtures.CreateAddressWithEmptyFields()
		result, err = common.MapCommonRequestAddress(*addressInvalid)
		assert.Error(t, err)
		assert.Nil(t, result)

		// Caso inválido: municipio no válido para el departamento
		addressInvalidMunicipality := fixtures.CreateAddressWithInvalidMunicipality()
		result, err = common.MapCommonRequestAddress(*addressInvalidMunicipality)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test para MapCommonRequestReceiver
	t.Run("TestMapCommonRequestReceiver", func(t *testing.T) {
		// Caso válido
		receiver := fixtures.CreateDefaultReceiver()
		result, err := common.MapCommonRequestReceiver(receiver)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, *receiver.DocumentType, result.DocumentType.GetValue())
		assert.Equal(t, *receiver.Name, *result.Name)

		// Caso inválido: email incorrecto
		receiverInvalidEmail := fixtures.CreateReceiverWithInvalidEmail()
		result, err = common.MapCommonRequestReceiver(receiverInvalidEmail)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test para MapCommonRequestExtension
	t.Run("TestMapCommonRequestExtension", func(t *testing.T) {
		// Caso válido
		extension := fixtures.CreateDefaultExtension()
		result, err := common.MapCommonRequestExtension(extension)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, extension.DeliveryName, result.DeliveryName.GetValue())
		assert.Equal(t, extension.ReceiverName, result.ReceiverName.GetValue())

		// Caso inválido: campos faltantes
		extensionInvalid := fixtures.CreateExtensionWithMissingFields()
		result, err = common.MapCommonRequestExtension(extensionInvalid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test para MapCommonRequestAppendix
	t.Run("TestMapCommonRequestAppendix", func(t *testing.T) {
		// Caso válido
		appendixes := []structs.AppendixRequest{fixtures.CreateDefaultAppendix()}
		result, err := common.MapCommonRequestAppendix(appendixes)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, appendixes[0].Field, result[0].Field.GetValue())

		// Caso inválido: campo vacío
		appendixesInvalid := []structs.AppendixRequest{fixtures.CreateAppendixWithInvalidField()}
		result, err = common.MapCommonRequestAppendix(appendixesInvalid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test para MapCommonRequestPaymentsType
	t.Run("TestMapCommonRequestPaymentsType", func(t *testing.T) {
		// Caso válido
		payments := []structs.PaymentRequest{fixtures.CreateDefaultPayment()}
		result, err := common.MapCommonRequestPaymentsType(payments)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, payments[0].Code, result[0].GetCode())

		// Caso inválido: código no válido
		paymentsInvalid := []structs.PaymentRequest{fixtures.CreatePaymentWithInvalidCode()}
		result, err = common.MapCommonRequestPaymentsType(paymentsInvalid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test para MapCommonRequestThirdPartySale
	t.Run("TestMapCommonRequestThirdPartySale", func(t *testing.T) {
		// Caso válido
		thirdPartySale := fixtures.CreateDefaultThirdPartySale()
		result, err := common.MapCommonRequestThirdPartySale(thirdPartySale)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, thirdPartySale.Name, result.Name)

		// Caso inválido: NIT vacío
		thirdPartySaleInvalid := fixtures.CreateThirdPartySaleWithEmptyNIT()
		result, err = common.MapCommonRequestThirdPartySale(thirdPartySaleInvalid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
