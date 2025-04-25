package dte

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	transmissionPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// GenericDTEUseCase implementa un caso de uso genérico para cualquier tipo de DTE
type GenericDTEUseCase struct {
	authService    auth.AuthManager
	dteService     transmissionPorts.DTEManager
	transmitter    appPorts.BaseTransmitter
	service        ports.DTEService
	mapper         mapper.DTEMapper
	responseMapper mapper.ResponseMapperFunc
	additionalOps  AdditionalOperationsFunc
}

// NewGenericDTEUseCase crea una nueva instancia de GenericDTEUseCase
func NewGenericDTEUseCase(
	authService auth.AuthManager,
	dteService transmissionPorts.DTEManager,
	transmitter appPorts.BaseTransmitter,
	service ports.DTEService,
	mapper mapper.DTEMapper,
	responseMapper mapper.ResponseMapperFunc,
	additionalOps AdditionalOperationsFunc,
) *GenericDTEUseCase {
	return &GenericDTEUseCase{
		authService:    authService,
		dteService:     dteService,
		transmitter:    transmitter,
		service:        service,
		mapper:         mapper,
		responseMapper: responseMapper,
		additionalOps:  additionalOps,
	}
}

// Create procesa cualquier tipo de DTE utilizando un flujo genérico
func (u *GenericDTEUseCase) Create(ctx context.Context, req interface{}) (interface{}, *response.SuccessOptions, error) {
	// 1. Obtener los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	token := ctx.Value("token").(string)

	// 2. Obtener la información del emisor
	issuer, err := u.authService.GetIssuer(ctx, claims.BranchID)
	if err != nil {
		logs.Error("Error getting issuer information", map[string]interface{}{"error": err.Error()})
		return nil, nil, err
	}

	// 3. Mapear a modelo de dominio
	domainModel, err := u.mapper.MapToDomainModel(req, issuer)
	if err != nil {
		logs.Error("Error mapping to domain model", map[string]interface{}{"error": err.Error()})
		return nil, nil, err
	}

	// 4. Crear DTE a nivel de servicio
	result, err := u.service.Create(ctx, domainModel, claims.BranchID)
	if err != nil {
		logs.Error("Error creating DTE at service level", map[string]interface{}{"error": err.Error()})
		return nil, nil, err
	}

	// 5. Mapear a modelo de hacienda
	mhModel := u.responseMapper(result)

	// 6. Extraer el código de generación
	generationCode, err := extractGenerationCode(mhModel)
	if err != nil {
		logs.Error("Error extracting generation code", map[string]interface{}{"error": err.Error()})
		return nil, nil, err
	}

	// 7. Configurar detalles de respuesta
	options := &response.SuccessOptions{
		Ambient:        config.Server.AmbientCode,
		GenerationCode: generationCode,
		EmissionDate:   utils.TimeNow(),
	}

	// 8. Comenzar la transmisión del documento
	transmitResult, err := u.transmitter.RetryTransmission(ctx, mhModel, token, claims.NIT)
	if err != nil {
		logs.Error("Error transmitting document", map[string]interface{}{"error": err.Error()})
		return mhModel, options, err
	}
	options.ReceptionStamp = transmitResult.ReceptionStamp

	// 9. Guardar el documento en la base de datos
	err = u.dteService.Create(ctx, mhModel, constants.TransmissionNormal, constants.DocumentReceived, transmitResult.ReceptionStamp)
	if err != nil {
		logs.Error("Error saving document in database", map[string]interface{}{"error": err.Error()})
		return mhModel, options, err
	}

	// 10. Ejecutar operaciones adicionales específicas (si las hay)
	if u.additionalOps != nil {
		err = u.additionalOps(ctx, result, claims.BranchID, mhModel)
		if err != nil {
			logs.Error("Error executing additional operations", map[string]interface{}{"error": err.Error()})
			return mhModel, options, err
		}
	}

	return mhModel, options, nil
}

// extractGenerationCode extrae el código de generación usando reflexión
func extractGenerationCode(mhModel interface{}) (string, error) {
	extractor, err := utils.ExtractAuxiliarIdentification(mhModel)
	if err != nil {
		return "", err
	}

	return extractor.Identification.GenerationCode, nil
}
