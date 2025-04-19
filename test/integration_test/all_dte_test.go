package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	transmitterModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/transmitter/hacienda_error"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/MarlonG1/api-facturacion-sv/test/fixtures"
	"github.com/MarlonG1/api-facturacion-sv/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DTETestConfig contiene la configuración para el tipo de documento a probar
type DTETestConfig struct {
	EndpointPath string
	DocumentType string
	RequestType  interface{}
	GetRequest   func() interface{}
	DteBuilder   func() (interfaces.DTEDocument, error)
	MapperConfig struct {
		RequestMapperAdapter mapper.DTEMapper
		ResponseMapper       mapper.ResponseMapperFunc
	}
}

func TestAllDTETypes(t *testing.T) {
	test.TestMain(t)

	factory := mapper.NewMapperFactory()

	// Configuraciones para cada tipo de documento
	dteConfigs := map[string]DTETestConfig{
		"Invoice": {
			EndpointPath: "/invoice",
			DocumentType: constants.FacturaElectronica,
			RequestType:  &structs.CreateInvoiceRequest{},
			GetRequest: func() interface{} {
				return fixtures.CreateDefaultInvoiceRequest()
			},
			DteBuilder: func() (interfaces.DTEDocument, error) {
				dteBuilder := fixtures.NewDTEBuilder()
				return dteBuilder.BuildElectronicInvoice()
			},
			MapperConfig: struct {
				RequestMapperAdapter mapper.DTEMapper
				ResponseMapper       mapper.ResponseMapperFunc
			}{
				RequestMapperAdapter: factory.CreateInvoiceMapperAdapter(),
				ResponseMapper:       factory.GetInvoiceResponseMapper(),
			},
		},
		"CCF": {
			EndpointPath: "/ccf",
			DocumentType: constants.CCFElectronico,
			RequestType:  &structs.CreateCreditFiscalRequest{},
			GetRequest: func() interface{} {
				return fixtures.CreateDefaultCreditFiscalRequest()
			},
			DteBuilder: func() (interfaces.DTEDocument, error) {
				dteBuilder := fixtures.NewDTEBuilder()
				return dteBuilder.BuildCreditFiscalDocument()
			},
			MapperConfig: struct {
				RequestMapperAdapter mapper.DTEMapper
				ResponseMapper       mapper.ResponseMapperFunc
			}{
				RequestMapperAdapter: factory.CreateCCFMapperAdapter(),
				ResponseMapper:       factory.GetCCFResponseMapper(),
			},
		},
		"CreditNote": {
			EndpointPath: "/creditnote",
			DocumentType: constants.NotaCreditoElectronica,
			RequestType:  &structs.CreateCreditNoteRequest{},
			GetRequest: func() interface{} {
				return fixtures.CreateDefaultCreditNoteRequest()
			},
			DteBuilder: func() (interfaces.DTEDocument, error) {
				dteBuilder := fixtures.NewDTEBuilder()
				return dteBuilder.BuildCreditNote()
			},
			MapperConfig: struct {
				RequestMapperAdapter mapper.DTEMapper
				ResponseMapper       mapper.ResponseMapperFunc
			}{
				RequestMapperAdapter: factory.CreateCreditNoteMapperAdapter(),
				ResponseMapper:       factory.GetCreditNoteResponseMapper(),
			},
		},
		"Retention": {
			EndpointPath: "/retention",
			DocumentType: constants.ComprobanteRetencionElectronico,
			RequestType:  &structs.CreateRetentionRequest{},
			GetRequest: func() interface{} {
				return fixtures.CreateMixedDocumentsRetentionRequest()
			},
			DteBuilder: func() (interfaces.DTEDocument, error) {
				dteBuilder := fixtures.NewDTEBuilder()
				return dteBuilder.BuildRetentionDocumentWithMixedItems()
			},
			MapperConfig: struct {
				RequestMapperAdapter mapper.DTEMapper
				ResponseMapper       mapper.ResponseMapperFunc
			}{
				RequestMapperAdapter: factory.CreateRetentionMapperAdapter(),
				ResponseMapper:       factory.GetRetentionResponseMapper(),
			},
		},
	}

	// Definir todos los escenarios de prueba
	testCases := []struct {
		name       string
		setupMocks func(mockAuthManager *mocks.MockAuthManager, mockDTEService *mocks.MockDTEService,
			mockDTEManager *mocks.MockDTEManager, mockTransmitter *mocks.MockBaseTransmitter, mockContingency *mocks.MockContingencyManager,
			dteConfig DTETestConfig)
		prepareRequest    func(dteConfig DTETestConfig) (*http.Request, error)
		expectedStatus    int
		validateResponse  func(t *testing.T, recorder *httptest.ResponseRecorder, dteConfig DTETestConfig)
		handleContingency bool
	}{
		{
			name: "Normal emission - success case",
			setupMocks: func(mockAuthManager *mocks.MockAuthManager, mockDTEService *mocks.MockDTEService,
				mockDTEManager *mocks.MockDTEManager, mockTransmitter *mocks.MockBaseTransmitter,
				mockContingency *mocks.MockContingencyManager, dteConfig DTETestConfig) {

				// 1. Mock para obtener el emisor
				issuer := fixtures.CreateDefaultIssuer()
				mockAuthManager.EXPECT().
					GetIssuer(gomock.Any(), uint(1)).
					Return(issuer, nil)

				// 2. Mock para crear a nivel de servicio
				mockDTE, err := dteConfig.DteBuilder()
				if err != nil {
					t.Fatalf("Error building DTE document: %v", err)
				}

				mockDTEService.EXPECT().
					Create(
						gomock.Any(),
						gomock.Any(),
						uint(1),
					).
					Return(mockDTE, nil)

				// 3. Mock para transmitir a Hacienda
				transmitResponse := &transmitterModels.TransmitResult{
					Status:         "PROCESADO",
					ReceptionStamp: utils.ToStringPointer("2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM"),
				}
				mockTransmitter.EXPECT().
					RetryTransmission(
						gomock.Any(),
						gomock.Any(),
						"test-token",
						"11111111111111",
					).
					Return(transmitResponse, nil)

				// 4. Mock para guardar el documento
				mockDTEManager.EXPECT().
					Create(
						gomock.Any(),
						gomock.Any(),
						constants.TransmissionNormal,
						constants.DocumentReceived,
						utils.ToStringPointer("2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM"),
					).
					Return(nil)

				// 5. Mock de contingencia - NUNCA DEBE SER LLAMADO EN ESTE CASO
				mockContingency.EXPECT().
					StoreDocumentInContingency(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).Times(0)
			},
			prepareRequest: func(dteConfig DTETestConfig) (*http.Request, error) {
				dteRequest := dteConfig.GetRequest()
				requestJSON, err := json.Marshal(dteRequest)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest("POST", "/api/v1"+dteConfig.EndpointPath, bytes.NewBuffer(requestJSON))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")

				// Agregar contexto con claims de autenticación simulados
				claims := &models.AuthClaims{
					BranchID: 1,
					NIT:      "11111111111111",
				}
				ctx := context.WithValue(req.Context(), "claims", claims)
				ctx = context.WithValue(ctx, "token", "test-token")
				req = req.WithContext(ctx)

				return req, nil
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, dteConfig DTETestConfig) {
				var response struct {
					Success        bool            `json:"success"`
					Data           json.RawMessage `json:"data"`
					ReceptionStamp string          `json:"reception_stamp"`
					QRLink         string          `json:"qr_link"`
				}

				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				// Verificar la respuesta general
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.ReceptionStamp)
				assert.NotEmpty(t, response.QRLink)
				assert.Equal(t, "2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM", response.ReceptionStamp)

				// Verificar identificación básica en todos los documentos
				var identification struct {
					Identificacion struct {
						Version       int    `json:"version"`
						Ambiente      string `json:"ambiente"`
						TipoDte       string `json:"tipoDte"`
						TipoOperacion int    `json:"tipoOperacion"`
						TipoModelo    int    `json:"tipoModelo"`
						NumeroControl string `json:"numeroControl"`
					} `json:"identificacion"`
				}

				err = json.Unmarshal(response.Data, &identification)
				require.NoError(t, err)

				assert.Equal(t, 1, identification.Identificacion.Version)
				assert.Equal(t, constants.Testing, identification.Identificacion.Ambiente)
				assert.Equal(t, dteConfig.DocumentType, identification.Identificacion.TipoDte)
				assert.Equal(t, constants.TransmisionNormal, identification.Identificacion.TipoOperacion)
				assert.Equal(t, constants.ModeloFacturacionPrevio, identification.Identificacion.TipoModelo)
				assert.NotEmpty(t, identification.Identificacion.NumeroControl)
			},
			handleContingency: false,
		},
		{
			name: "Contingency emission - success case",
			setupMocks: func(mockAuthManager *mocks.MockAuthManager, mockDTEService *mocks.MockDTEService,
				mockDTEManager *mocks.MockDTEManager, mockTransmitter *mocks.MockBaseTransmitter,
				mockContingency *mocks.MockContingencyManager, dteConfig DTETestConfig) {

				// 1. Mock para obtener el emisor
				issuer := fixtures.CreateDefaultIssuer()
				mockAuthManager.EXPECT().
					GetIssuer(gomock.Any(), uint(1)).
					Return(issuer, nil)

				// 2. Mock para crear a nivel de servicio
				mockDTE, err := dteConfig.DteBuilder()
				if err != nil {
					t.Fatalf("Error building DTE document: %v", err)
				}

				mockDTEService.EXPECT().
					Create(
						gomock.Any(),
						gomock.Any(),
						uint(1),
					).
					Return(mockDTE, nil)

				// 3. En contingencia no hay transmisión exitosa
				mockTransmitter.EXPECT().
					RetryTransmission(
						gomock.Any(),
						gomock.Any(),
						"test-token",
						"11111111111111",
					).
					Return(nil, &hacienda_error.HTTPResponseError{
						StatusCode: http.StatusServiceUnavailable,
						Body:       []byte("Forced contingency - service unavailable"),
						URL:        config.MHPaths.ReceptionURL,
						Method:     "POST",
					}).
					AnyTimes()

				// 4. Mock de contingencia - En caso de error, se guarda en contingencia
				mockContingency.EXPECT().
					StoreDocumentInContingency(
						gomock.Any(),
						gomock.Any(),
						dteConfig.DocumentType,
						int8(constants.NoDisponibilidadMH),
						constants.ContingencyReasons[constants.NoDisponibilidadMH],
					).Return(nil)
			},
			prepareRequest: func(dteConfig DTETestConfig) (*http.Request, error) {
				dteRequest := dteConfig.GetRequest()
				requestJSON, err := json.Marshal(dteRequest)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest("POST", "/api/v1"+dteConfig.EndpointPath, bytes.NewBuffer(requestJSON))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")

				// Agregar contexto con claims de autenticación simulados
				claims := &models.AuthClaims{
					BranchID: 1,
					NIT:      "11111111111111",
				}
				ctx := context.WithValue(req.Context(), "claims", claims)
				ctx = context.WithValue(ctx, "token", "test-token")
				req = req.WithContext(ctx)

				return req, nil
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, dteConfig DTETestConfig) {
				var response struct {
					Success        bool            `json:"success"`
					Data           json.RawMessage `json:"data"`
					ReceptionStamp *string         `json:"reception_stamp"`
					QRLink         string          `json:"qr_link"`
				}

				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				// Verificar la respuesta general
				assert.True(t, response.Success)
				assert.Nil(t, response.ReceptionStamp) // Sin sello en contingencia

				// Verificar identificación con campos de contingencia
				var identification struct {
					Identificacion struct {
						Version          int    `json:"version"`
						Ambiente         string `json:"ambiente"`
						TipoDte          string `json:"tipoDte"`
						TipoOperacion    int    `json:"tipoOperacion"`
						TipoModelo       int    `json:"tipoModelo"`
						NumeroControl    string `json:"numeroControl"`
						TipoContingencia int    `json:"tipoContingencia"`
						MotivoContin     string `json:"motivoContin"`
					} `json:"identificacion"`
				}

				err = json.Unmarshal(response.Data, &identification)
				require.NoError(t, err)

				assert.Equal(t, dteConfig.DocumentType, identification.Identificacion.TipoDte)
				assert.Equal(t, constants.TransmisionContingencia, identification.Identificacion.TipoOperacion)
				assert.Equal(t, constants.ModeloFacturacionDiferido, identification.Identificacion.TipoModelo)
				assert.Equal(t, constants.NoDisponibilidadMH, identification.Identificacion.TipoContingencia)
				assert.Equal(t, constants.ContingencyReasons[constants.NoDisponibilidadMH], identification.Identificacion.MotivoContin)
			},
			handleContingency: true,
		},
		{
			name: "Validation error - error case",
			setupMocks: func(mockAuthManager *mocks.MockAuthManager, mockDTEService *mocks.MockDTEService,
				mockDTEManager *mocks.MockDTEManager, mockTransmitter *mocks.MockBaseTransmitter,
				mockContingency *mocks.MockContingencyManager, dteConfig DTETestConfig) {

				// 1. Mock para obtener el emisor
				issuer := fixtures.CreateDefaultIssuer()
				mockAuthManager.EXPECT().
					GetIssuer(gomock.Any(), uint(1)).
					Return(issuer, nil)

				// 2. Mock para crear a nivel de servicio - devuelve error de validación
				mockDTEService.EXPECT().
					Create(
						gomock.Any(),
						gomock.Any(),
						uint(1),
					).
					Return(nil, dte_errors.NewValidationError("RequiredField", "Request->Receiver"))

				// 3. Mock de contingencia - NUNCA DEBE SER LLAMADO EN ESTE CASO
				mockContingency.EXPECT().
					StoreDocumentInContingency(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).Times(0)
			},
			prepareRequest: func(dteConfig DTETestConfig) (*http.Request, error) {
				// Obtener request por defecto y modificarlo para que sea inválido
				dteRequest := dteConfig.GetRequest()
				doInvalidAmount(dteRequest)

				requestJSON, err := json.Marshal(dteRequest)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest("POST", "/api/v1"+dteConfig.EndpointPath, bytes.NewBuffer(requestJSON))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")

				// Agregar contexto con claims de autenticación simulados
				claims := &models.AuthClaims{
					BranchID: 1,
					NIT:      "11111111111111",
				}
				ctx := context.WithValue(req.Context(), "claims", claims)
				ctx = context.WithValue(ctx, "token", "test-token")
				req = req.WithContext(ctx)

				return req, nil
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, dteConfig DTETestConfig) {
				var response struct {
					Success bool            `json:"success"`
					Error   json.RawMessage `json:"error"`
				}

				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				// Verificar que hay un error
				assert.False(t, response.Success)
				assert.NotNil(t, response.Error)

				// Verificar que el error contiene la palabra "required"
				errorStr := string(response.Error)
				assert.Contains(t, errorStr, "required")
			},
			handleContingency: false,
		},
	}

	// Para cada tipo de DTE, ejecutar todas las pruebas
	for dteName, dteConfig := range dteConfigs {
		t.Run(dteName, func(t *testing.T) {
			// Para cada caso de prueba
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					// Crear los mocks
					mockAuthManager := mocks.NewMockAuthManager(ctrl)
					mockDTEService := mocks.NewMockDTEService(ctrl)
					mockDTEManager := mocks.NewMockDTEManager(ctrl)
					mockTransmitter := mocks.NewMockBaseTransmitter(ctrl)
					mockContingency := mocks.NewMockContingencyManager(ctrl)

					// Configurar los mocks según el caso de prueba
					tc.setupMocks(mockAuthManager, mockDTEService, mockDTEManager, mockTransmitter, mockContingency, dteConfig)

					// Crear el caso de uso para el tipo de documento
					var additionalOps dte.AdditionalOperationsFunc = nil
					genericUseCase := dte.NewGenericDTEUseCase(
						mockAuthManager,
						mockDTEManager,
						mockTransmitter,
						mockDTEService,
						dteConfig.MapperConfig.RequestMapperAdapter,
						dteConfig.MapperConfig.ResponseMapper,
						additionalOps,
					)

					// Configurar el handler
					contingencyHandler := helpers.NewContingencyHandler(mockContingency)
					genericHandler := handlers.NewGenericDTEHandler(contingencyHandler)

					// Registrar el documento
					genericHandler.RegisterDocument(dteConfig.EndpointPath, helpers.DocumentConfig{
						DocumentType:    dteConfig.DocumentType,
						UseCase:         genericUseCase,
						RequestType:     dteConfig.RequestType,
						UsesContingency: true,
					})

					// Preparar la solicitud HTTP para este tipo de DTE
					req, err := tc.prepareRequest(dteConfig)
					require.NoError(t, err)

					// Ejecutar la solicitud
					recorder := httptest.NewRecorder()
					router := mux.NewRouter()
					router.HandleFunc("/api/v1"+dteConfig.EndpointPath, genericHandler.HandleCreate).Methods("POST")
					router.ServeHTTP(recorder, req)

					// Verificar el código de estado HTTP
					assert.Equal(t, tc.expectedStatus, recorder.Code)

					// Verificar la respuesta según el caso de prueba
					tc.validateResponse(t, recorder, dteConfig)
				})
			}
		})
	}
}

// doInvalidAmount modifica el monto de un DTE para que sea inválido a nivel de dominio
func doInvalidAmount(
	dteRequest interface{},
) {
	invalidValue := 40.0
	switch req := dteRequest.(type) {
	case *structs.CreateInvoiceRequest:
		req.Items[0].TaxedSale = invalidValue
	case *structs.CreateCreditFiscalRequest:
		req.Items[0].TaxedSale = invalidValue
	case *structs.CreateCreditNoteRequest:
		req.Items[0].TaxedSale = invalidValue
	case *structs.CreateRetentionRequest:
		req.Items[0].TaxedAmount = &invalidValue
	}
}
