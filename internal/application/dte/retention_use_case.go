package dte

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	transmissionPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	requestDTO "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type RetentionUseCase struct {
	dteManager       transmissionPorts.DTEManager
	authService      auth.AuthManager
	retentionService retention.RetentionManager
	transmitter      appPorts.BaseTransmitter
	mapper           *request_mapper.RetentionMapper
}

func NewRetentionUseCase(authService auth.AuthManager, retentionService retention.RetentionManager, dteService transmissionPorts.DTEManager, transmitter appPorts.BaseTransmitter) *RetentionUseCase {
	return &RetentionUseCase{
		authService:      authService,
		dteManager:       dteService,
		retentionService: retentionService,
		transmitter:      transmitter,
		mapper:           request_mapper.NewRetentionMapper(),
	}
}

func (u *RetentionUseCase) Create(ctx context.Context, req *requestDTO.CreateRetentionRequest) (*structs.RetentionDTEResponse, *response.SuccessOptions, error) {
	// 1. Obtener los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	token := ctx.Value("token").(string)

	// 2. Obtener la información del emisor
	issuer, err := u.authService.GetIssuer(ctx, claims.BranchID)
	if err != nil {
		return nil, nil, err
	}

	// 3. Mapear a modelo de dominio
	inputRetention, err := u.mapper.MapToRetentionData(req, issuer)
	if err != nil {
		return nil, nil, err
	}

	// 4. Crear retencion a nivel de servicio
	retention, err := u.retentionService.Create(ctx, inputRetention, claims.BranchID, u.mapper.IsAllPhysical(req.Items))
	if err != nil {
		return nil, nil, err
	}

	// 5. Mapear a modelo de hacienda
	retentionMh, err := response_mapper.ToRetentionMH(retention)
	if err != nil {
		return nil, nil, err
	}
	// 6. Mapear a modelo de respuesta
	options := &response.SuccessOptions{
		Ambient:        config.Server.AmbientCode,
		GenerationCode: retentionMh.Identificacion.CodigoGeneracion,
		EmissionDate:   utils.TimeNow(),
	}

	// 7. Comenzar la transmisión de la factura
	result, err := u.transmitter.RetryTransmission(ctx, retentionMh, token, claims.NIT)
	if err != nil {
		return retentionMh, options, err
	}
	options.ReceptionStamp = result.ReceptionStamp

	if result.Status != ReceivedStatus {
		logs.Warn("Error transmitting invoice", map[string]interface{}{"error": "TransmissionFailed"})
		return retentionMh, options, dte_errors.NewDTEErrorSimple("TransmissionFailed")
	}

	// 8. Guardar la factura en la base de datos
	err = u.dteManager.Create(ctx, retentionMh, constants.TransmissionNormal, constants.DocumentReceived, result.ReceptionStamp)
	if err != nil {
		return retentionMh, options, err
	}

	return retentionMh, options, nil
}
