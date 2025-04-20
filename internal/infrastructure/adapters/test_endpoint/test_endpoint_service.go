package test_endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	commonModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	identificationVO "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/location"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type testService struct {
	db *gorm.DB
}

func NewTestService(db *gorm.DB) test_endpoint.TestManager {
	return &testService{
		db: db,
	}
}

func (s *testService) RunSystemTest() (*models.TestResult, error) {
	startTime := utils.TimeNow()
	tests := make([]models.ComponentTest, 0)

	// Test 1: Base de datos
	dbTest := s.testDatabase()
	tests = append(tests, dbTest)
	if !dbTest.Success {
		return s.buildResult(tests, startTime), nil
	}

	// Test 2: Mapeo y validación de DTE
	mappingTest := s.testDTEMapping()
	tests = append(tests, mappingTest)
	if !mappingTest.Success {
		return s.buildResult(tests, startTime), nil
	}

	// Test 3: Generación de números de control
	seqTest := s.testSequenceGeneration()
	tests = append(tests, seqTest)

	// Test 4: Transmisión a Hacienda
	testDTE := getTestDTE()
	haciendaTest := s.testHaciendaTransmission(testDTE)
	tests = append(tests, haciendaTest)

	return s.buildResult(tests, startTime), nil
}

func (s *testService) testDatabase() models.ComponentTest {
	start := utils.TimeNow()
	test := models.ComponentTest{
		Name: "database_connection",
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		logs.Error("Database connection test failed", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
	}

	if err := sqlDB.Ping(); err != nil {
		logs.Error("Database ping test failed", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
	} else {
		test.Success = true
	}

	test.Duration = time.Since(start).Milliseconds()
	return test
}

func (s *testService) testDTEMapping() models.ComponentTest {
	start := utils.TimeNow()
	test := models.ComponentTest{
		Name: "dte_mapping",
	}

	// DTE de prueba predefinido
	testDTE := getTestDTE()

	mh := response_mapper.ToMHInvoice(testDTE)
	if mh == nil {
		logs.Error("DTE mapping test failed")
		test.Success = false
	} else {
		test.Success = true
	}

	test.Duration = time.Since(start).Milliseconds()
	return test
}

func (s *testService) testSequenceGeneration() models.ComponentTest {
	start := utils.TimeNow()
	test := models.ComponentTest{
		Name: "sequence_generation",
	}

	// Intenta generar un número de secuencia de prueba
	err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&db_models.ControlNumberSequence{}).
			Where("branch_id = ? AND dte_type = ?", 0, "01").
			UpdateColumn("last_number", gorm.Expr("last_number + ?", 1)).
			Error
	})

	if err != nil {
		logs.Error("Sequence generation test failed", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
	} else {
		test.Success = true
	}

	test.Duration = time.Since(start).Milliseconds()
	return test
}

func (s *testService) buildResult(tests []models.ComponentTest, startTime time.Time) *models.TestResult {
	success := true
	for _, test := range tests {
		if !test.Success {
			success = false
			break
		}
	}

	return &models.TestResult{
		Success:  success,
		Tests:    tests,
		Duration: time.Since(startTime).Milliseconds(),
	}
}

func (s *testService) testHaciendaTransmission(testDTE *invoice_models.ElectronicInvoice) models.ComponentTest {
	start := utils.TimeNow()
	test := models.ComponentTest{
		Name: "hacienda_transmission",
	}

	mhDTE := response_mapper.ToMHInvoice(testDTE)
	if mhDTE == nil {
		logs.Error("Failed to map test DTE")
		test.Success = false
		test.Duration = time.Since(start).Milliseconds()
		return test
	}

	// Preparar request a Hacienda
	jsonData, err := json.Marshal(mhDTE)
	if err != nil {
		logs.Error("Failed to marshal test DTE", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
		test.Duration = time.Since(start).Milliseconds()
		return test
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		"POST",
		config.MHPaths.ReceptionURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		logs.Error("Failed to create request", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
		test.Duration = time.Since(start).Milliseconds()
		return test
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error("Failed to transmit to Hacienda", map[string]interface{}{
			"error": err.Error(),
		})
		test.Success = false
		test.Duration = time.Since(start).Milliseconds()
		return test
	}
	defer resp.Body.Close()

	// Se considera exitosa la prueba si:
	// 1. Recibimos una respuesta (no hubo error de conexión)
	// 2. El código de estado está entre 400-499 (error esperado al ser DTE de prueba)
	test.Success = resp.StatusCode >= 400 && resp.StatusCode < 500

	logs.Info("Hacienda transmission test completed", map[string]interface{}{
		"statusCode": resp.StatusCode,
		"success":    test.Success,
	})

	test.Duration = time.Since(start).Milliseconds()
	return test
}

// getTestDTE genera un DTE con el mínimo requerido para pruebas
func getTestDTE() *invoice_models.ElectronicInvoice {
	identification, err := common.MapCommonRequestIdentification(1, 1, constants.FacturaElectronica)
	if err != nil {
		logs.Error("Failed to map identification", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	issuer := &commonModels.Issuer{
		NIT:                 *identificationVO.NewValidatedNIT("TEST PARA PRUEBAS"),
		NRC:                 *identificationVO.NewValidatedNRC("TEST PARA PRUEBAS"),
		Name:                "EMPRESA DE PRUEBA",
		ActivityCode:        *identificationVO.NewValidatedActivityCode("01234"),
		ActivityDescription: "ACTIVIDAD DE PRUEBA",
		EstablishmentType:   *document.NewValidatedEstablishmentType("01"),
		Address: &commonModels.Address{
			Department:   *location.NewValidatedDepartment("01"),
			Municipality: *location.NewValidatedMunicipality("01", "01"),
			Complement:   *location.NewValidatedAddress("DIRECCION DE PRUEBA"),
		},
		Phone: *base.NewValidatedPhone("22222222"),
		Email: *base.NewValidatedEmail("test@test.com"),
	}

	name := "CLIENTE DE PRUEBA"
	receiver := &commonModels.Receiver{
		DocumentType:   document.NewValidatedDTEType("13"),
		DocumentNumber: identificationVO.NewValidatedDocumentNumber("00000000-0"),
		Name:           &name,
		Address: &commonModels.Address{
			Department:   *location.NewValidatedDepartment("09"),
			Municipality: *location.NewValidatedMunicipality("04", "09"),
			Complement:   *location.NewValidatedAddress("SIMON"),
		},
		Email: base.NewValidatedEmail("example@gmail.com"),
	}

	items := []invoice_models.InvoiceItem{
		{
			Item: &commonModels.Item{
				Number:      *item.NewValidatedItemNumber(1),
				Type:        *item.NewValidatedItemType(1),
				Description: "PRODUCTO DE PRUEBA",
				Quantity:    *item.NewValidatedQuantity(1),
				UnitMeasure: *item.NewValidatedUnitMeasure(59),
				UnitPrice:   *financial.NewValidatedAmount(7.50),
				Discount:    *financial.NewValidatedDiscount(0.525),
				Code:        item.NewValidatedItemCode("6609"),
				Taxes:       []string{constants.TaxIVA},
			},
			NonSubjectSale: *financial.NewValidatedAmount(0),
			ExemptSale:     *financial.NewValidatedAmount(0),
			TaxedSale:      *financial.NewValidatedAmount(6.98),
			SuggestedPrice: *financial.NewValidatedAmount(0),
			NonTaxed:       *financial.NewValidatedAmount(0),
			IVAItem:        *financial.NewValidatedAmount(0.80),
		},
	}

	summary := invoice_models.InvoiceSummary{
		Summary: &commonModels.Summary{
			TotalNonSubject:    *financial.NewValidatedAmount(0),
			TotalExempt:        *financial.NewValidatedAmount(0),
			TotalTaxed:         *financial.NewValidatedAmount(6.98),
			SubTotal:           *financial.NewValidatedAmount(6.98),
			SubTotalSales:      *financial.NewValidatedAmount(6.98),
			NonSubjectDiscount: *financial.NewValidatedAmount(0),
			ExemptDiscount:     *financial.NewValidatedAmount(0),
			DiscountPercentage: *financial.NewValidatedDiscount(0),
			TotalDiscount:      *financial.NewValidatedAmount(0.53),
			TotalOperation:     *financial.NewValidatedAmount(6.98),
			TotalNonTaxed:      *financial.NewValidatedAmount(0),
			OperationCondition: *financial.NewValidatedPaymentCondition(1),
			TotalToPay:         *financial.NewValidatedAmount(6.98),
			TotalTaxes: []interfaces.Tax{
				&commonModels.Tax{
					Code:        *financial.NewValidatedTaxType(constants.TaxIVA),
					Description: "IVA 13%",
					Value:       &commonModels.TaxAmount{TotalAmount: *financial.NewValidatedAmount(0.80)},
				},
			},
			PaymentTypes: []interfaces.PaymentType{
				&commonModels.PaymentType{
					Code:      *financial.NewValidatedPaymentType("01"),
					Amount:    *financial.NewValidatedAmount(6.98),
					Reference: "",
				},
			},
		},
		TotalIva: *financial.NewValidatedAmount(0.80),
	}

	var itemsInterface = make([]interfaces.Item, len(items))
	for i, item := range items {
		itemsInterface[i] = &item
	}

	return &invoice_models.ElectronicInvoice{
		DTEDocument: &commonModels.DTEDocument{
			Identification: identification,
			Issuer:         issuer,
			Receiver:       receiver,
			Items:          itemsInterface,
			Summary:        &summary,
		},
		InvoiceItems:   items,
		InvoiceSummary: summary,
	}
}
