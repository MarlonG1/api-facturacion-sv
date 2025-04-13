package mappers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

func TestMapToRetentionData(t *testing.T) {
	// 1. Inicializar la configuración del entorno
	rootPath := utils.FindProjectRoot()
	err := config.InitEnvConfig(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Inicializar el tiempo global
	err = utils.TimeInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Emisor de ejemplo:
	issuer := &dte.IssuerDTE{
		NIT:                  "11111111111111",
		NRC:                  "1111111",
		CommercialName:       "JEMEPLO",
		BusinessName:         "EMPRESA DE PRUEBAS SA DE CV",
		EconomicActivity:     "11111",
		EconomicActivityDesc: "Venta al por mayor de otros productos",
		EstablishmentCode:    utils.ToStringPointer("C001"),
		Email:                utils.ToStringPointer("email@gmail.com"),
		Phone:                utils.ToStringPointer("22567890"),
		Address: &user.Address{
			Department:   "06",
			Municipality: "20",
			Complement:   "BOULEVARD SANTA ELENA SUR, SANTA TECLA",
		},
		EstablishmentType:   "02",
		EstablishmentCodeMH: nil,
		POSCode:             nil,
		POSCodeMH:           nil,
	}

	// Casos de prueba
	tests := []struct {
		name          string
		reqJSON       string
		wantErr       bool
		errorContains string
		checkReceiver bool
	}{
		{
			name: "All physical documents",
			reqJSON: `{
				"items": [
					{
						"type": 1,
						"document_number": "S221001345",
						"description": "Compra de suministros de oficina",
						"retention_code": "22",
						"taxed_amount": 115.25,
						"iva_amount": 19.00,
						"emission_date": "2025-03-20",
						"dte_type": "03"
					},
					{
						"type": 1,
						"document_number": "S221001346",
						"description": "Servicio de limpieza",
						"retention_code": "C4",
						"taxed_amount": 226.50,
						"iva_amount": 19.00,
						"emission_date": "2025-03-22",
						"dte_type": "03"
					}
				],
				"summary": {
					"total_retention_amount": 341.75,
					"total_retention_iva": 44.43
				},
				"receiver": {
					"document_type": "36",
					"document_number": "06141804941035",
					"nrc": "123456",
					"name": "Empresa Servicios Generales, S.A. de C.V.",
					"commercial_name": "ServiGeneral",
					"activity_code": "46900",
					"activity_description": "Venta al por mayor de otros productos",
					"address": {
						"department": "06",
						"municipality": "20",
						"complement": "Colonia Escalón, Calle La Reforma #123, San Salvador"
					},
					"phone": "22123456",
					"email": "google@gmail.com"
				}
			}`,
			wantErr:       false,
			checkReceiver: true,
		},
		{
			name: "All physical documents without optional fields",
			reqJSON: `{
				"items": [
					{
						"type": 1,
						"document_number": "S221001345",
						"description": "Compra de suministros de oficina",
						"retention_code": "22",
						"taxed_amount": 115.25,
						"iva_amount": 19.00,
						"emission_date": "2025-03-20",
						"dte_type": "03"
					},
					{
						"type": 1,
						"document_number": "S221001346",
						"description": "Servicio de limpieza",
						"retention_code": "C4",
						"taxed_amount": 226.50,
						"iva_amount": 19.00,
						"emission_date": "2025-03-22",
						"dte_type": "03"
					}
				],
				"summary": {
					"total_retention_amount": 341.75,
					"total_retention_iva": 44.43
				},
				"receiver": {
					"document_type": "36",
					"document_number": "06141804941035",
					"name": "Empresa Servicios Generales, S.A. de C.V.",
					"activity_code": "46900",
					"activity_description": "Venta al por mayor de otros productos",
					"address": {
						"department": "06",
						"municipality": "20",
						"complement": "Colonia Escalón, Calle La Reforma #123, San Salvador"
					},
					"email": "google@gmail.com"
				}
			}`,
			wantErr:       false,
			checkReceiver: true,
		},
		{
			name: "All electronic documents",
			reqJSON: `{
				"items": [
					{
						"type": 2,
						"document_number": "FF54E9DB-79C3-42CE-B432-EC522C97EFB9",
						"description": "Compra de equipos informáticos",
						"retention_code": "22"
					},
					{
						"type": 2,
						"document_number": "AD54E9BB-79A3-42AE-B432-EC522C97EFB7",
						"description": "Mantenimiento de servidores",
						"retention_code": "C4"
					}
				],
				"receiver": {
					"document_type": "36",
					"document_number": "06141804941035",
					"nrc": "123456",
					"name": "Empresa Servicios Generales, S.A. de C.V.",
					"commercial_name": "ServiGeneral",
					"activity_code": "46900",
					"activity_description": "Venta al por mayor de otros productos",
					"address": {
						"department": "06",
						"municipality": "20",
						"complement": "Colonia Escalón, Calle La Reforma #123, San Salvador"
					},
					"phone": "22123456",
					"email": "google@gmail.com"
				},
				"extension": {
					"observation": "Retención por servicios tecnológicos primer trimestre",
					"delivery_name": "Juan Carlos Martínez",
					"delivery_document": "04567890-1",
					"receiver_name": "Ana María López",
					"receiver_document": "12345678-9"
				}
			}`,
			wantErr:       false,
			checkReceiver: false,
		},
		{
			name: "Mixed documents (physical and electronic)",
			reqJSON: `{
				"items": [
					{
						"type": 1,
						"document_number": "S221001347",
						"description": "Consultoría financiera",
						"retention_code": "C9",
						"taxed_amount": 450.00,
						"iva_amount": 19.00,
						"emission_date": "2025-03-15",
						"dte_type": "03"
					},
					{
						"type": 2,
						"document_number": "FF32E9DB-79C3-42CE-B432-EC522C97EFB2",
						"description": "Servicios de auditoría",
						"retention_code": "C4"
					}
				],
				"receiver": {
					"document_type": "36",
					"document_number": "06141804941035",
					"nrc": "123456",
					"name": "Empresa Servicios Generales, S.A. de C.V.",
					"commercial_name": "ServiGeneral",
					"activity_code": "46900",
					"activity_description": "Venta al por mayor de otros productos",
					"address": {
						"department": "06",
						"municipality": "20",
						"complement": "Colonia Escalón, Calle La Reforma #123, San Salvador"
					},
					"phone": "22123456",
					"email": "google@gmail.com"
				},
				"appendixes": [
					{
						"field": "nota_interna",
						"label": "Nota interna",
						"value": "Sisisisisisisisisisisisisisisiisis"
					}
				]
			}`,
			wantErr:       false,
			checkReceiver: false,
		},
		{
			name: "Error: Missing required field (receiver) for physical document",
			reqJSON: `{
				"items": [
					{
						"type": 1,
						"document_number": "S221001348",
						"description": "Servicio de consultoría",
						"retention_code": "22",
						"emission_date": "2025-03-20",
						"dte_type": "03"
					}
				]
			}`,
			wantErr:       true,
			errorContains: "RequiredField",
		},
		{
			name: "Error: Missing required field (receiver) for electronic document",
			reqJSON: `{
				"items": [
					{
						"type": 2,
						"document_number": "1EEAB582-AA75-4D9C-AFFE-47E7FF89D24E",
						"description": "Servicio de consultoría",
						"retention_code": "22",
						"emission_date": "2025-03-20",
						"dte_type": "03"
					}
				]
			}`,
			wantErr:       true,
			errorContains: "RequiredField",
		},
		{
			name: "Error: Invalid retention code for electronic document",
			reqJSON: `{
				"items": [
					{
						"type": 2,
						"document_number": "FF54E9DB-79C3-42CE-B432-EC522C97EFB9",
						"description": "Compra de equipos informáticos",
						"retention_code": "99"
					}
				]
			}`,
			wantErr:       true,
			errorContains: "It must be one of the allowed retention codes",
		},
	}

	mapper := request_mapper.NewRetentionMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req structs.CreateRetentionRequest
			err := json.Unmarshal([]byte(tt.reqJSON), &req)
			assert.NoError(t, err, "Error unmarshal JSON")

			// Llamar a la función a probar
			got, err := mapper.MapToRetentionData(&req, issuer)

			// Verificar resultados
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotNil(t, got.RetentionItems)
			assert.Len(t, got.RetentionItems, len(req.Items))

			// Verificar campos específicos según el tipo de prueba
			if tt.checkReceiver {
				assert.NotNil(t, got.InputDataCommon.Receiver, "The receiver should be present")
				assert.NotNil(t, got.RetentionSummary, "The summary should be present")

				for i, item := range got.RetentionItems {
					assert.NotNil(t, item, "The item at index %v", i)
					assert.NotNil(t, item.RetentionAmount, "The item at index %v", i)
					assert.NotNil(t, item.RetentionIVA, "The item at index %v", i)
				}
			}

			// Verificar campos opcionales
			if req.Extension != nil {
				assert.NotNil(t, got.Extension)
			}

			if req.Appendixes != nil {
				assert.NotNil(t, got.Appendixes)
				assert.Len(t, got.Appendixes, len(req.Appendixes))
			}
		})
	}
}
