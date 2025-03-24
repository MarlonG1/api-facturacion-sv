package dte

import (
	"context"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	transmissionInterface "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	requestDTO "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

type CCFUseCase struct {
	authService ports.AuthManager
	ccfService  interfaces.CCFManager
	dteService  transmissionInterface.DTEManager
	transmitter appPorts.BaseTransmitter
	mapper      *request_mapper.CCFMapper
}

func NewCCFUseCase(authService ports.AuthManager, invoiceService interfaces.CCFManager, transmitter appPorts.BaseTransmitter, dteService transmissionInterface.DTEManager) *CCFUseCase {
	return &CCFUseCase{
		authService: authService,
		ccfService:  invoiceService,
		transmitter: transmitter,
		dteService:  dteService,
		mapper:      request_mapper.NewCCFMapper(),
	}
}

func (u *CCFUseCase) Create(ctx context.Context, req *requestDTO.CreateCreditFiscalRequest) (*structs.CCFDTEResponse, *string, error) {
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

	// 6. Comenzar la transmisión de la factura
	result, err := u.transmitter.RetryTransmission(ctx, mhCCF, token, claims.NIT)
	if err != nil {
		return mhCCF, nil, err
	}
	if result.Status != ReceivedStatus {
		return mhCCF, result.ReceptionStamp, dte_errors.NewDTEErrorSimple("TransmissionFailed")
	}

	// 7. Guardar el CCF en la base de datos
	err = u.dteService.Create(ctx, mhCCF, constants.TransmissionNormal, constants.DocumentReceived, result.ReceptionStamp)
	if err != nil {
		return mhCCF, result.ReceptionStamp, err
	}

	return mhCCF, result.ReceptionStamp, nil
}
