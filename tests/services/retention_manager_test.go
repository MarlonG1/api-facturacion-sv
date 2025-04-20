package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/MarlonG1/api-facturacion-sv/tests"
	"github.com/MarlonG1/api-facturacion-sv/tests/fixtures"
	"github.com/MarlonG1/api-facturacion-sv/tests/mocks"
)

func TestRetentionServiceCreate(t *testing.T) {
	test.TestMain(t)

	tests := []struct {
		name               string
		setupRetentionData func() (*retention_models.InputRetentionData, error)
		setupMock          func(*mocks.MockSequentialNumberManager, *mocks.MockDTEManager)
		wantErr            bool
		errorCode          string
	}{
		{
			name: "Valid Retention with physical items",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildValidRetention()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockSeq.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.ComprobanteRetencionElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-06-R0010001-000000000012345", nil)
			},
			wantErr: false,
		},
		{
			name: "Valid Retention with electronic items",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithElectronicItems()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockSeq.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.ComprobanteRetencionElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-06-R0010001-000000000012345", nil)

				// Mock para cada documento electrónico
				mockDTE.EXPECT().GetByGenerationCode(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(&dte.DTEDocument{
					CreatedAt: utils.TimeNow(),
					Details: &dte.DTEDetails{
						DTEType:  constants.FacturaElectronica,
						JSONData: `{"resumen":{"subTotal":100.00,"ivaRete1":13.00}}`,
					},
				}, nil).AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "Valid Retention with mixed items",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithMixedItems()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockSeq.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.ComprobanteRetencionElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-06-R0010001-000000000012345", nil)

				// Mock para documentos electrónicos
				mockDTE.EXPECT().GetByGenerationCode(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(&dte.DTEDocument{
					CreatedAt: utils.TimeNow(),
					Details: &dte.DTEDetails{
						DTEType:  constants.FacturaElectronica,
						JSONData: `{"resumen":{"subTotal":100.00,"ivaRete1":13.00}}`,
					},
				}, nil).AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "Error - Retention with inconsistent summary totals",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildInvalidRetentionWithInconsistentSummary()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				// No se espera ninguna llamada porque la validación fallará antes
			},
			wantErr:   true,
			errorCode: "InvalidTotalSubjectRetention",
		},
		{
			name: "Error - Retention with physical items but no summary",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithoutSummary()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				// No se espera ninguna llamada porque la validación fallará antes
			},
			wantErr:   true,
			errorCode: "InvalidTotalSubjectRetention",
		},
		{
			name: "Error - Retention with invalid IVA calculation",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildValidRetention()
				if err != nil {
					return nil, err
				}

				// Modificar la retención de IVA para que sea incorrecta
				retention.RetentionItems[0].RetentionIVA = *financial.NewValidatedAmount(50.0) // Valor incorrecto

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				// No se espera ninguna llamada porque la validación fallará antes
			},
			wantErr:   true,
			errorCode: "InvalidRetentionIVA",
		},
		{
			name: "Error - Retention with document date out of range",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildValidRetention()
				if err != nil {
					return nil, err
				}

				// Establecer una fecha fuera del período permitido (más de un mes atrás)
				oldDate, _ := temporal.NewEmissionDate(time.Now().AddDate(0, -3, 0))
				retention.RetentionItems[0].EmissionDate = *oldDate

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				// No se espera ninguna llamada porque la validación fallará antes
			},
			wantErr:   true,
			errorCode: "DateOutOfAllowedRange",
		},
		{
			name: "Error - Retention with invalid document types",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithElectronicItems()
				if err != nil {
					return nil, err
				}

				inputData := fixtures.BuildAsInputRetentionData(retention)

				// Asegurar que el tipo de documento se reconoce como electrónico
				for i := range inputData.RetentionItems {
					docType, _ := document.NewOperationType(constants.ElectronicDocument)
					inputData.RetentionItems[i].DocumentType = *docType
				}

				return inputData, nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockDTE.EXPECT().GetByGenerationCode(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(&dte.DTEDocument{
					CreatedAt: utils.TimeNow(),
					Details: &dte.DTEDetails{
						DTEType:  constants.DocumentInvalid, // Tipo inválido
						JSONData: `{"resumen":{"subTotal":100.00,"ivaRete1":13.00}}`,
					},
				}, nil).AnyTimes()
			},
			wantErr:   true,
			errorCode: "InvalidDTETypeForRetention",
		},
		{
			name: "Error - Retention with DTE extraction error",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithElectronicItems()
				if err != nil {
					return nil, err
				}

				inputData := fixtures.BuildAsInputRetentionData(retention)

				// Asegurar que el tipo de documento se reconoce como electrónico
				for i := range inputData.RetentionItems {
					docType, _ := document.NewOperationType(constants.ElectronicDocument)
					inputData.RetentionItems[i].DocumentType = *docType
				}

				return inputData, nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockDTE.EXPECT().GetByGenerationCode(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(&dte.DTEDocument{
					CreatedAt: utils.TimeNow(),
					Details: &dte.DTEDetails{
						DTEType:  constants.FacturaElectronica,
						JSONData: `{"invalid_json": "invalid"}}`, // JSON inválido para provocar error en extracción
					},
				}, nil).AnyTimes()
			},
			wantErr:   true,
			errorCode: "invalid character",
		},
		{
			name: "Error - Retention with document not found",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildRetentionWithElectronicItems()
				if err != nil {
					return nil, err
				}

				inputData := fixtures.BuildAsInputRetentionData(retention)

				// Asegurar que se reconoce como electrónico
				for i := range inputData.RetentionItems {
					docType, _ := document.NewOperationType(constants.ElectronicDocument)
					inputData.RetentionItems[i].DocumentType = *docType
				}

				return inputData, nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockDTE.EXPECT().GetByGenerationCode(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, fmt.Errorf("document not found")).AnyTimes()
			},
			wantErr:   true,
			errorCode: "document not found",
		},
		{
			name: "Error - Failed to generate control number",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildValidRetention()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockSeq.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.ComprobanteRetencionElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("", shared_error.NewGeneralServiceError("SequentialNumberManager", "GetNextControlNumber", "Failed to generate control number", nil))
			},
			wantErr:   true,
			errorCode: "Failed to generate control number",
		},
		{
			name: "Error - Invalid control number format",
			setupRetentionData: func() (*retention_models.InputRetentionData, error) {
				retention, err := fixtures.BuildValidRetention()
				if err != nil {
					return nil, err
				}

				return fixtures.BuildAsInputRetentionData(retention), nil
			},
			setupMock: func(mockSeq *mocks.MockSequentialNumberManager, mockDTE *mocks.MockDTEManager) {
				mockSeq.EXPECT().GetNextControlNumber(
					gomock.Any(),
					constants.ComprobanteRetencionElectronico,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return("DTE-06-R001-INVALID", nil) // Formato inválido
			},
			wantErr:   true,
			errorCode: "InvalidPattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retentionData, err := tt.setupRetentionData()
			if err != nil {
				t.Fatalf("Error preparing test data: %v", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSeqNumberManager := mocks.NewMockSequentialNumberManager(ctrl)
			mockDTEManager := mocks.NewMockDTEManager(ctrl)
			tt.setupMock(mockSeqNumberManager, mockDTEManager)

			service := retention.NewRetentionService(mockSeqNumberManager, mockDTEManager)

			result, err := service.Create(context.Background(), retentionData, 1)

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

				retentionDoc, ok := result.(*retention_models.RetentionModel)
				assert.True(t, ok, "Result should be a RetentionModel")
				assert.NotNil(t, retentionDoc.Identification)
				assert.NotNil(t, retentionDoc.Issuer)
				assert.NotNil(t, retentionDoc.Receiver)
				assert.NotEmpty(t, retentionDoc.RetentionItems)

				if len(retentionDoc.RetentionItems) > 0 && retentionDoc.RetentionItems[0].DocumentType.GetValue() == constants.PhysicalDocument {
					assert.NotNil(t, retentionDoc.RetentionSummary)
					assert.True(t, retentionDoc.RetentionSummary.TotalSubjectRetention.GetValue() > 0)
					assert.True(t, retentionDoc.RetentionSummary.TotalIVARetention.GetValue() > 0)
				}

				assert.NotEmpty(t, retentionDoc.Identification.GetControlNumber())
				assert.NotEmpty(t, retentionDoc.Identification.GetGenerationCode())
			}
		})
	}
}
