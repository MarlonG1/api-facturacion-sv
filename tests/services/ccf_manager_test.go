package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/MarlonG1/api-facturacion-sv/tests"
	"github.com/MarlonG1/api-facturacion-sv/tests/fixtures"
	"github.com/MarlonG1/api-facturacion-sv/tests/mocks"
)

func TestCreditFiscalServiceCreate(t *testing.T) {
	test.TestMain(t)

	tests := []struct {
		name         string
		setupCCFData func() (*ccf_models.CCFData, error)
		setupMock    func(*mocks.MockSequentialNumberManager)
		wantErr      bool
		errorCode    string
	}{
		{
			name: "Valid CCF creation",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with IVA perception (1%)",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer una percepción válida (1% del total gravado)
				taxedAmount := ccf.CreditSummary.Summary.GetTotalTaxed()
				validPerception := financial.NewValidatedAmount(math.Round(taxedAmount*0.01*100) / 100)
				ccf.CreditSummary.IVAPerception = *validPerception

				// Actualizar el total a pagar para incluir la percepción
				totalOperation := ccf.CreditSummary.Summary.GetTotalOperation()
				newTotalToPay := totalOperation + validPerception.GetValue()
				ccf.CreditSummary.Summary.SetTotalToPay(newTotalToPay)

				// Actualizar los pagos
				payments := ccf.CreditSummary.Summary.GetPaymentTypes()
				for i, payment := range payments {
					if i == 0 { // Ajustar el primer pago
						payment.SetAmount(newTotalToPay)
					}
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with exempt sales only",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Configurar todas las ventas como exentas
				for i := range ccf.CreditItems {
					amount := ccf.CreditItems[i].TaxedSale.GetValue()
					taxedZero := financial.NewValidatedAmount(0.0)
					exemptAmount := financial.NewValidatedAmount(amount)

					ccf.CreditItems[i].TaxedSale = *taxedZero
					ccf.CreditItems[i].ExemptSale = *exemptAmount
					ccf.CreditItems[i].Taxes = nil // Sin impuestos para ventas exentas
				}

				// Actualizar los totales en el resumen
				totalExempt := 0.0
				for _, item := range ccf.CreditItems {
					totalExempt += item.ExemptSale.GetValue()
				}

				ccf.CreditSummary.Summary.SetTotalTaxed(0.0)
				ccf.CreditSummary.Summary.SetTotalExempt(totalExempt)
				ccf.CreditSummary.Summary.SetSubtotalSales(totalExempt)
				ccf.CreditSummary.Summary.SetSubTotal(totalExempt)
				ccf.CreditSummary.Summary.SetTotalOperation(totalExempt)
				ccf.CreditSummary.Summary.SetTotalToPay(totalExempt)
				ccf.CreditSummary.Summary.SetTotalTaxes(nil) // Sin impuestos

				// Actualizar los pagos
				payments := ccf.CreditSummary.Summary.GetPaymentTypes()
				for i, payment := range payments {
					if i == 0 { // Ajustar el primer pago
						payment.SetAmount(totalExempt)
					}
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with IVA and income retentions",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer retenciones válidas
				taxedAmount := ccf.CreditSummary.Summary.GetTotalTaxed()
				ivaRetention := financial.NewValidatedAmount(taxedAmount * 0.01)    // 1%
				incomeRetention := financial.NewValidatedAmount(taxedAmount * 0.05) // 5%

				ccf.CreditSummary.IVARetention = *ivaRetention
				ccf.CreditSummary.IncomeRetention = *incomeRetention

				// Actualizar el total a pagar para incluir las retenciones
				totalOperation := ccf.CreditSummary.Summary.GetTotalOperation()
				newTotalToPay := totalOperation - ivaRetention.GetValue() - incomeRetention.GetValue()
				newTotalToPay = math.Round(newTotalToPay*100) / 100
				ccf.CreditSummary.Summary.SetTotalToPay(newTotalToPay)

				// Actualizar los pagos
				payments := ccf.CreditSummary.Summary.GetPaymentTypes()
				for i, payment := range payments {
					if i == 0 {
						payment.SetAmount(newTotalToPay)
					}
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with related documents",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Añadir documentos relacionados válidos
				relatedDoc1 := models.RelatedDocument{}
				relatedDoc1.SetGenerationType(constants.ElectronicDocument)
				relatedDoc1.SetDocumentNumber("001BEDAD-93F3-4F49-85D9-1E3618425F6B")
				relatedDoc1.SetDocumentType(constants.NotaRemisionElectronica)
				relatedDoc1.SetEmissionDate(utils.TimeNow())

				relatedDoc2 := models.RelatedDocument{}
				relatedDoc2.SetGenerationType(constants.ElectronicDocument)
				relatedDoc2.SetDocumentNumber("001BEDAD-93F3-4F49-85D9-1E3618425F6B")
				relatedDoc2.SetDocumentType(constants.NotaRemisionElectronica)
				relatedDoc2.SetEmissionDate(utils.TimeNow())

				ccfData := fixtures.BuildAsCCFData(ccf)
				ccfData.RelatedDocs = []models.RelatedDocument{relatedDoc1, relatedDoc2}

				// Establecer la referencia a documentos relacionados en los ítems
				for i := range ccfData.Items {
					docRef := "001BEDAD-93F3-4F49-85D9-1E3618425F6B"
					ccfData.Items[i].SetRelatedDoc(&docRef)
				}

				return ccfData, nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with tax item (type 4)",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Añadir un ítem de tipo impuesto (4) válido
				taxItem := ccf.CreditItems[0]
				taxItem.Item.SetNumber(3) // Nuevo número de ítem
				taxItem.Item.SetType(constants.Impuesto)
				taxItem.Item.SetUnitMeasure(99) // Unidad de medida requerida para tipo 4
				taxItem.Item.SetDescription("IVA")

				// Actualizar los taxes de summary
				totalTaxedSummary := ccf.CreditSummary.Summary.GetTotalTaxed()
				totalTaxedSummary += taxItem.TaxedSale.GetValue()
				ccf.CreditSummary.Summary.SetTotalTaxed(totalTaxedSummary)
				ccf.CreditSummary.SetSubtotalSales(totalTaxedSummary)
				ccf.CreditSummary.Summary.SetSubTotal(totalTaxedSummary)

				// Agregar el item a los existentes
				ccf.CreditItems = append(ccf.CreditItems, taxItem)

				// Los impuestos en resumen
				taxes := models.Tax{}
				taxes.SetCode(constants.TaxIVA)
				taxes.SetValue(totalTaxedSummary * 0.13)
				taxes.SetDescription("IVA 13%")

				// Añadir el nuevo ítem a los existentes
				var taxInterfaces []interfaces.Tax
				taxInterfaces = append(taxInterfaces, &taxes)
				ccf.Summary.SetTotalTaxes(taxInterfaces)

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid CCF with maximum allowed related documents (50)",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Añadir exactamente 50 documentos relacionados (límite máximo)
				relatedDocs := make([]models.RelatedDocument, 50)
				for i := 0; i < 50; i++ {
					doc := models.RelatedDocument{}
					docNum := fmt.Sprintf("001BEDAD-93F3-4F49-85D9-1E3618425F%02d", i)
					doc.SetDocumentNumber(docNum)
					doc.SetGenerationType(constants.ElectronicDocument)
					doc.SetDocumentType(constants.DocContableLiquidacionElectronico)
					doc.SetEmissionDate(utils.TimeNow())
					relatedDocs[i] = doc
				}

				ccfData := fixtures.BuildAsCCFData(ccf)
				ccfData.RelatedDocs = relatedDocs

				// Establecer la referencia a un documento relacionado en los ítems
				for i := range ccfData.Items {
					docRef := "001BEDAD-93F3-4F49-85D9-1E3618425F00"
					ccfData.Items[i].SetRelatedDoc(&docRef)
				}

				return ccfData, nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Error - CCF with InvalidMixedSalesWithExempt",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Crear un ítem con ventas exentas y gravadas a la vez
				exemptAmount, _ := financial.NewAmount(100.0)
				ccf.CreditItems[0].ExemptSale = *exemptAmount

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidMixedSalesWithExempt",
		},
		{
			name: "Error - CCF with InvalidTaxesWithExempt",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer ventas exentas pero mantener impuestos
				taxedAmount, _ := financial.NewAmount(0.0)
				exemptAmount, _ := financial.NewAmount(100.0)
				ccf.CreditItems[0].TaxedSale = *taxedAmount
				ccf.CreditItems[0].ExemptSale = *exemptAmount
				// Mantener la lista de impuestos del ítem

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidTaxesWithExempt",
		},
		{
			name: "Error - CCF with InvalidMixedSalesWithNonSubject",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Crear un ítem con ventas no sujetas y gravadas a la vez
				nonSubjectAmount, _ := financial.NewAmount(100.0)
				ccf.CreditItems[0].NonSubjectSale = *nonSubjectAmount

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidMixedSalesWithNonSubject",
		},
		{
			name: "Error - CCF with InvalidMixedSalesWithNonTaxed",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Crear un ítem con montos no gravados y ventas gravadas a la vez
				nonTaxedAmount, _ := financial.NewAmount(100.0)
				ccf.CreditItems[0].NonTaxed = *nonTaxedAmount

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidMixedSalesWithNonTaxed",
		},
		{
			name: "Error - CCF with InvalidUnitPriceWithNonTaxed",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Crear un ítem con monto no gravado y precio unitario no cero
				nonTaxedAmount, _ := financial.NewAmount(100.0)
				zeroAmount, _ := financial.NewAmount(0.0)

				ccf.CreditItems[0].NonTaxed = *nonTaxedAmount
				ccf.CreditItems[0].TaxedSale = *zeroAmount
				ccf.CreditItems[0].ExemptSale = *zeroAmount
				ccf.CreditItems[0].NonSubjectSale = *zeroAmount
				ccf.CreditItems[0].Taxes = nil
				ccf.CreditItems[0].Item.SetForceUnitPrice(50.0) // Precio unitario no debe ser cero

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidUnitPriceWithNonTaxed",
		},
		{
			name: "Error - CCF with InvalidUnitMeasure",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Crear un ítem de tipo impuesto con unidad de medida incorrecta
				ccf.CreditItems[0].Item.SetType(constants.Impuesto)
				ccf.CreditItems[0].Item.SetUnitMeasure(58) // No es 99

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidUnitMeasure",
		},
		{
			name: "Error - CCF with DiscountExceedsSubtotal",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer un descuento que excede el subtotal
				subtotal := ccf.CreditSummary.Summary.GetSubTotal()
				ccf.CreditSummary.TaxedDiscount = *financial.NewValidatedAmount(subtotal * 2) // Descuento mayor al subtotal

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "DiscountExceedsSubtotal",
		},
		{
			name: "Error - CCF with MissingTaxes",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Mantener venta gravada pero eliminar los impuestos
				ccf.CreditSummary.Summary.SetTotalTaxes(nil)

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "MissingTaxes",
		},
		{
			name: "Error - CCF with InvalidPerceptionAmount",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer una percepción con valor incorrecto
				// La percepción debería ser el 1% de la venta gravada
				taxedAmount := ccf.CreditSummary.Summary.GetTotalTaxed()
				incorrectPerception := financial.NewValidatedAmount(taxedAmount * 0.02) // Debería ser 0.01
				ccf.CreditSummary.IVAPerception = *incorrectPerception

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidPerceptionAmount",
		},
		{
			name: "Error - CCF with RequiredField missing",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Eliminar el código de actividad económica del receptor
				ccf.Receiver.SetActivityCode(nil)

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Error - CCF with InvalidSubTotalCalculation",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Cambiar el subtotal a un valor incorrecto
				incorrectSubtotal := ccf.CreditSummary.Summary.GetSubtotalSales() -
					ccf.CreditSummary.TaxedDiscount.GetValue() -
					ccf.CreditSummary.Summary.GetExemptDiscount() -
					ccf.CreditSummary.Summary.GetNonSubjectDiscount() + 50.0 // Valor incorrecto

				ccf.CreditSummary.Summary.SetSubTotal(incorrectSubtotal)

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidSubTotalCalculation",
		},
		{
			name: "Error - CCF with InvalidMonetaryAmount",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Establecer un total a pagar con más de 2 decimales
				ccf.CreditSummary.Summary.SetForceTotalToPay(100.123) // 3 decimales

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidMonetaryAmount",
		},
		{
			name: "Error - CCF with ExceededRelatedDocsLimit",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Añadir más de 50 documentos relacionados
				relatedDocs := make([]models.RelatedDocument, 51)
				for i := 0; i < 51; i++ {
					doc := models.RelatedDocument{}
					docNum := fmt.Sprintf("DTE-01-C0020000-00000000000000%d", i)
					doc.SetDocumentNumber(docNum)
					doc.SetDocumentType(constants.FacturaElectronica)
					doc.SetEmissionDate(utils.TimeNow())
					relatedDocs[i] = doc
				}

				ccfData := fixtures.BuildAsCCFData(ccf)
				ccfData.RelatedDocs = relatedDocs

				return ccfData, nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "ExceededRelatedDocsLimit",
		},
		{
			name: "Error - CCF with InvalidRelatedDocDTEType",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				// Añadir un documento relacionado con un tipo inválido
				relatedDoc := models.RelatedDocument{}
				relatedDoc.SetDocumentNumber("DTE-01-C0020000-000000000000001")
				relatedDoc.SetDocumentType(constants.DocumentInvalid) // Tipo inválido

				relatedDoc.SetEmissionDate(utils.TimeNow())

				ccfData := fixtures.BuildAsCCFData(ccf)
				ccfData.RelatedDocs = []models.RelatedDocument{relatedDoc}

				return ccfData, nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidRelatedDocDTEType",
		},
		{
			name: "Error - CCF without NRC",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildInvalidCCFWithoutNRC()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "MissingNRC",
		},
		{
			name: "Error - CCF with invalid tax calculation",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				for _, t := range ccf.CreditSummary.Summary.GetTotalTaxes() {
					tax, ok := t.(*models.Tax)
					if ok && tax.GetCode() == constants.TaxIVA {
						tax.SetValue(50.0) // Incorrect value
						break
					}
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidTaxCalculation",
		},
		{
			name: "Error - CCF with invalid total to pay",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				ccf.CreditSummary.Summary.SetTotalToPay(1500.0) // Incorrect value

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {},
			wantErr:   true,
			errorCode: "InvalidTotalToPayCalculation",
		},
		{
			name: "Error - Failed to generate control number",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("", shared_error.NewGeneralServiceError("SequentialNumberManager", "GetNextControlNumber", "Failed to generate control number", nil))
			},
			wantErr:   true,
			errorCode: "Failed to generate control number",
		},
		{
			name: "Error - Failed to set control number",
			setupCCFData: func() (*ccf_models.CCFData, error) {
				ccf, err := fixtures.BuildCCF()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsCCFData(ccf), nil
			},
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C001000001-000000000012345", nil)
			},
			wantErr:   true,
			errorCode: "InvalidPattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccfData, err := tt.setupCCFData()
			if err != nil {
				t.Fatalf("Error preparing test data: %v", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSeqNumberManager := mocks.NewMockSequentialNumberManager(ctrl)
			tt.setupMock(mockSeqNumberManager)

			service := ccf.NewCCFService(mockSeqNumberManager)

			result, err := service.Create(context.Background(), ccfData, 1)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorCode != "" {
					var dteErr *dte_errors.DTEError
					var serviceErr *shared_error.ServiceError

					if errors.As(err, &dteErr) {
						assert.Contains(t, dteErr.Error(), tt.errorCode, "Error message should contain expected code")
					} else if errors.As(err, &serviceErr) {
						assert.Contains(t, serviceErr.Error(), tt.errorCode, "Error message should contain expected code")
					} else {
						assert.Contains(t, err.Error(), tt.errorCode, "Error message should contain expected code")
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				ccfDoc, ok := result.(*ccf_models.CreditFiscalDocument)
				assert.True(t, ok, "Result should be a CreditFiscalDocument")
				assert.NotNil(t, ccfDoc.Identification)
				assert.NotNil(t, ccfDoc.Issuer)
				assert.NotNil(t, ccfDoc.Receiver)
				assert.NotEmpty(t, ccfDoc.CreditItems)
				assert.NotNil(t, ccfDoc.CreditSummary.Summary)

				assert.NotEmpty(t, ccfDoc.Identification.GetControlNumber())
				assert.NotEmpty(t, ccfDoc.Identification.GetGenerationCode())
			}
		})
	}
}

func TestCreditFiscalServiceCreateWithDifferentOperationConditions(t *testing.T) {
	test.TestMain(t)

	tests := []struct {
		name               string
		operationCondition int
		setupMock          func(*mocks.MockSequentialNumberManager)
		wantErr            bool
	}{
		{
			name:               "Valid CCF with cash operation",
			operationCondition: constants.Cash,
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name:               "Valid CCF with credit operation",
			operationCondition: constants.Credit,
			setupMock: func(mock *mocks.MockSequentialNumberManager) {
				mock.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.CCFElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-03-C0010001-000000000012345", nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := fixtures.NewCCFBuilder()

			var ccfModel *ccf_models.CreditFiscalDocument
			var err error

			if tt.operationCondition == constants.Credit {
				builder.AddIdentification().
					AddIssuer().
					AddReceiverForCCF().
					AddItems().
					AddSummaryWithCreditOperation()

				ccfModel, err = builder.BuildWithoutValidation()
				if err != nil {
					t.Fatalf("Error building CCF: %v", err)
				}
			} else {
				ccfModel, err = fixtures.BuildCCF()
				if err != nil {
					t.Fatalf("Error building CCF: %v", err)
				}
			}

			ccfData := fixtures.BuildAsCCFData(ccfModel)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSeqNumberManager := mocks.NewMockSequentialNumberManager(ctrl)
			tt.setupMock(mockSeqNumberManager)

			service := ccf.NewCCFService(mockSeqNumberManager)

			result, err := service.Create(context.Background(), ccfData, 1)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				ccfDoc, ok := result.(*ccf_models.CreditFiscalDocument)
				assert.True(t, ok)
				assert.Equal(t, tt.operationCondition, ccfDoc.Summary.GetOperationCondition())
			}
		})
	}
}

func TestCreditFiscalServiceCreateWithDifferentItems(t *testing.T) {
	test.TestMain(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSeqNumberManager := mocks.NewMockSequentialNumberManager(ctrl)
	mockSeqNumberManager.EXPECT().GetNextControlNumber(
		gomock.Any(),
		constants.CCFElectronico,
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return("DTE-03-C0010001-000000000012345", nil).AnyTimes()

	service := ccf.NewCCFService(mockSeqNumberManager)

	ccf, err := fixtures.BuildCCFWithMixedItemsType()
	if err != nil {
		t.Fatalf("Error building CCF: %v", err)
	}

	ccfData := fixtures.BuildAsCCFData(ccf)
	result, err := service.Create(context.Background(), ccfData, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	ccfDoc, ok := result.(*ccf_models.CreditFiscalDocument)
	assert.True(t, ok)

	foundProduct := false
	foundService := false

	for _, item := range ccfDoc.CreditItems {
		if item.GetType() == constants.Producto {
			foundProduct = true
		} else if item.GetType() == constants.Servicio {
			foundService = true
		}
	}

	assert.True(t, foundProduct, "Should have at least one Product type item")
	assert.True(t, foundService, "Should have at least one Service type item")
}
