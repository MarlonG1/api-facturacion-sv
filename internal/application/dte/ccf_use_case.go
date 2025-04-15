package dte

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	transmissionInterface "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	requestDTO "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type CCFUseCase struct {
	authService auth.AuthManager
	ccfService  ccf.CCFManager
	dteService  transmissionInterface.DTEManager
	transmitter appPorts.BaseTransmitter
	mapper      *request_mapper.CCFMapper
}

func NewCCFUseCase(authService auth.AuthManager, invoiceService ccf.CCFManager, transmitter appPorts.BaseTransmitter, dteService transmissionInterface.DTEManager) *CCFUseCase {
	return &CCFUseCase{
		authService: authService,
		ccfService:  invoiceService,
		transmitter: transmitter,
		dteService:  dteService,
		mapper:      request_mapper.NewCCFMapper(),
	}
}

func (u *CCFUseCase) Create(ctx context.Context, req *requestDTO.CreateCreditFiscalRequest) (*structs.CCFDTEResponse, *response.SuccessOptions, error) {
	// 1. Obtener los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	token := ctx.Value("token").(string)

	// 2. Obtener la información del emisor
	issuer, err := u.authService.GetIssuer(ctx, claims.BranchID)
	if err != nil {
		return nil, nil, err
	}

	// 3. Mapear a modelo de dominio
	reqCCF, err := u.mapper.MapToCCFData(req, issuer)
	if err != nil {
		return nil, nil, err
	}

	// 4. Crear invoice electrónica a nivel de servicio
	ccf, err := u.ccfService.Create(ctx, reqCCF, claims.BranchID)
	if err != nil {
		return nil, nil, err
	}

	// 5. Mapear a modelo de hacienda
	mhCCF, err := response_mapper.ToMHCreditFiscalInvoice(ccf)
	if err != nil {
		return nil, nil, err
	}

	// 6. Configurar detalles de respuesta
	options := &response.SuccessOptions{
		Ambient:        config.Server.AmbientCode,
		GenerationCode: mhCCF.Identificacion.CodigoGeneracion,
		EmissionDate:   utils.TimeNow(),
	}

	// 7. Comenzar la transmisión de la factura
	result, err := u.transmitter.RetryTransmission(ctx, mhCCF, token, claims.NIT)
	if err != nil {
		return mhCCF, options, err
	}
	options.ReceptionStamp = result.ReceptionStamp

	if result.Status != ReceivedStatus {
		return mhCCF, options, dte_errors.NewDTEErrorSimple("TransmissionFailed")
	}

	// 8. Guardar el CCF en la base de datos
	err = u.dteService.Create(ctx, mhCCF, constants.TransmissionNormal, constants.DocumentReceived, result.ReceptionStamp)
	if err != nil {
		return mhCCF, options, err
	}

	return mhCCF, options, nil
}
