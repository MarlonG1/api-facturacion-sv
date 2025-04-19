package mappers

import (
	"testing"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/MarlonG1/api-facturacion-sv/test"
	"github.com/MarlonG1/api-facturacion-sv/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMapToInvalidationData(t *testing.T) {
	test.TestMain(t)

	// Emisor por defecto para todas las pruebas
	issuer := fixtures.CreateDefaultIssuer()

	// Crear documentos base para pruebas
	invoiceDTE := createInvoiceDTE()
	ccfDTE := createCCFDTE()

	// Definir casos de prueba
	tests := []struct {
		name      string
		req       func() *structs.CreateInvalidationRequest
		baseDTE   *dte.DTEDetails
		wantErr   bool
		errorCode string
	}{
		// ------ VALIDACIONES BÁSICAS ------
		{
			name: "Valid invalidation request for invoice with replacement",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateDefaultInvalidationRequest()
			},
			baseDTE: invoiceDTE,
			wantErr: false,
		},
		{
			name: "Valid invalidation request for CCF with replacement",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateDefaultInvalidationRequest()
			},
			baseDTE: ccfDTE,
			wantErr: false,
		},
		{
			name: "Null invalidation request",
			req: func() *structs.CreateInvalidationRequest {
				return nil
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation with null base DTE",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateDefaultInvalidationRequest()
			},
			baseDTE:   nil,
			wantErr:   true,
			errorCode: "ErrorMapping",
		},

		// ------ VALIDACIONES DE TIPOS DE INVALIDACIÓN ------
		{
			name: "Invalidation with invalid type",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateInvalidationWithInvalidType()
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationType",
		},
		{
			name: "Type 3 invalidation without reason",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateInvalidationTypeWithoutReason()
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationType3",
		},
		{
			name: "Type 2 invalidation with replacement code",
			req: func() *structs.CreateInvalidationRequest {
				return fixtures.CreateInvalidationType2WithReplacementCode()
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationType2",
		},
		{
			name: "Valid type 2 invalidation (annulment)",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.Type = 2                 // Anulación
				req.ReplacementGenerationCode = nil // Sin código de reemplazo
				return req
			},
			baseDTE: invoiceDTE,
			wantErr: false,
		},
		{
			name: "Valid type 3 invalidation (definitive)",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.Type = 3 // Invalidación definitiva
				reason := "Documento con errores graves no recuperables"
				req.Reason.Reason = &reason
				return req
			},
			baseDTE: invoiceDTE,
			wantErr: false,
		},

		// ------ VALIDACIONES DE DOCUMENTOS RESPONSABLES ------
		{
			name: "Invalidation without responsible name",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.ResponsibleName = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation without responsible document type",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.ResponsibleDocType = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation without responsible document number",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.ResponsibleNumDoc = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation with invalid responsible document type",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.ResponsibleDocType = "99" // Tipo inválido
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidDocumentForReceiver",
		},

		// Nota, la validación que si el número de documento es un formato correcto para DUI o NIT se hace a nivel de dominio,
		// por lo tanto no se valida en el mapper

		// ------ VALIDACIONES DE DOCUMENTOS SOLICITANTES ------
		{
			name: "Invalidation without requestor name",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.RequestorName = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation without requestor document type",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.RequestorDocType = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation without requestor document number",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.RequestorNumDoc = ""
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invalidation with invalid requestor document type",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.RequestorDocType = "99" // Tipo inválido
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidDocumentForReceiver",
		},

		// Nota, la validación que si el número de documento es un formato correcto para DUI o NIT se hace a nivel de dominio,
		// por lo tanto no se valida en el mapper

		// ------ VALIDACIONES DE RAZÓN DE INVALIDACIÓN ------
		{
			name: "Type 3 invalidation with empty reason",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.Type = 3 // Invalidación definitiva
				emptyReason := ""
				req.Reason.Reason = &emptyReason
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationReason",
		},
		{
			name: "Type 3 invalidation with reason too short",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.Type = 3   // Invalidación definitiva
				shortReason := "abcd" // Menos de 5 caracteres
				req.Reason.Reason = &shortReason
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationReason",
		},
		{
			name: "Type 3 invalidation with reason too long",
			req: func() *structs.CreateInvalidationRequest {
				req := fixtures.CreateDefaultInvalidationRequest()
				req.Reason.Type = 3                          // Invalidación definitiva
				tooLongReason := createStringWithLength(251) // Más de 250 caracteres
				req.Reason.Reason = &tooLongReason
				return req
			},
			baseDTE:   invoiceDTE,
			wantErr:   true,
			errorCode: "InvalidInvalidationReason",
		},
	}

	// Ejecutar casos de prueba
	mapper := request_mapper.NewInvalidationMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req()

			// Validar primero la solicitud
			if req != nil {
				validationErr := mapper.ValidateInvalidationReRequest(req)
				if validationErr != nil {
					// Si hay error de validación y esperamos error, verificamos y retornamos
					if tt.wantErr {
						assert.Error(t, validationErr)
						if tt.errorCode != "" {
							assertErrorCode(t, validationErr, tt.errorCode)
						}
						return
					}
					// Si hay error pero no esperamos error, fallamos la prueba
					t.Fatalf("Unexpected validation error: %v", validationErr)
				}
			}

			// Proceder con el mapeo
			emissionDate := time.Now()
			got, err := mapper.MapToInvalidationData(req, issuer, tt.baseDTE, emissionDate)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorCode != "" {
					assertErrorCode(t, err, tt.errorCode)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)

			// Verificar correspondencia con el documento base
			if tt.baseDTE != nil {
				assert.Equal(t, tt.baseDTE.DTEType, got.Document.Type.GetValue())
				assert.Equal(t, tt.baseDTE.ID, got.Document.GenerationCode.GetValue())
				assert.Equal(t, tt.baseDTE.ControlNumber, got.Document.ControlNumber.GetValue())
			}

			// Verificar campos requeridos
			assert.NotNil(t, got.Identification)
			assert.NotNil(t, got.Reason)
			assert.NotNil(t, got.Issuer)
			assert.NotNil(t, got.Document)

			// Verificar tipo de invalidación
			assert.Equal(t, req.Reason.Type, int(got.Reason.Type.GetValue()))

			// Verificar código de reemplazo si aplica
			if req.ReplacementGenerationCode != nil {
				assert.NotNil(t, got.Document.ReplacementCode)
				assert.Equal(t, *req.ReplacementGenerationCode, got.Document.ReplacementCode.GetValue())
			} else {
				assert.Nil(t, got.Document.ReplacementCode)
			}

			// Verificar razón para tipo 3
			if req.Reason.Type == 3 {
				assert.NotNil(t, got.Reason.Reason)
				assert.Equal(t, *req.Reason.Reason, got.Reason.Reason.GetValue())
			}
		})
	}
}

// Funciones auxiliares para crear DTEs de prueba
func createInvoiceDTE() *dte.DTEDetails {
	return &dte.DTEDetails{
		ID:             "FF54E9DB-79C3-42CE-B432-EC522C97EFB9",
		DTEType:        "01", // Factura Electrónica
		ControlNumber:  "DTE-01-00000000-000000000000001",
		ReceptionStamp: utils.ToStringPointer("2025AAFEEE1A566A44F19A622C0C35C8A1B6FAZM"),
		JSONData: `{
			"receptor": {
				"nombre": "Cliente Ejemplo",
				"telefono": "22123456",
				"correo": "cliente@example.com",
				"tipoDocumento": "13",
				"numDocumento": "01234567-8"
			},
			"resumen": {
				"totalIva": 13.00
			}
		}`,
	}
}

func createCCFDTE() *dte.DTEDetails {
	return &dte.DTEDetails{
		ID:             "AD54E9BB-79A3-42AE-B432-EC522C97EFB7",
		DTEType:        "03", // CCF Electrónico
		ControlNumber:  "DTE-03-00000000-000000000000001",
		ReceptionStamp: utils.ToStringPointer("2025BBFEEE1A566A44F19A622C0C35C8A1B6FAZM"),
		JSONData: `{
			"receptor": {
				"nombre": "Empresa Cliente, S.A. de C.V.",
				"telefono": "22123456",
				"correo": "empresa@example.com",
				"nit": "06141804941035",
				"nrc": "1234567"
			},
			"resumen": {
				"totalIva": 13.00
			}
		}`,
	}
}

// Función auxiliar para crear una cadena de cierta longitud
func createStringWithLength(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s += "a"
	}
	return s
}
