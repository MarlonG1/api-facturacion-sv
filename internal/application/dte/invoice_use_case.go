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
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	requestDTO "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type InvoiceUseCase struct {
	authService    auth.AuthManager
	invoiceService invoice.InvoiceManager
	dteService     transmissionPorts.DTEManager
	transmitter    appPorts.BaseTransmitter
	mapper         *request_mapper.InvoiceMapper
}

func NewInvoiceUseCase(authService auth.AuthManager, invoiceService invoice.InvoiceManager,
	transmitter appPorts.BaseTransmitter, dteService transmissionPorts.DTEManager) *InvoiceUseCase {
	return &InvoiceUseCase{
		authService:    authService,
		invoiceService: invoiceService,
		transmitter:    transmitter,
		dteService:     dteService,
		mapper:         request_mapper.NewInvoiceMapper(),
	}
}

func (u *InvoiceUseCase) Create(ctx context.Context, req *requestDTO.CreateInvoiceRequest) (*structs.InvoiceDTEResponse, *response.SuccessOptions, error) {
	// 1. Obtener los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	token := ctx.Value("token").(string)

	// 2. Obtener la información del emisor
	issuer, err := u.authService.GetIssuer(ctx, claims.BranchID)
	if err != nil {
		return nil, nil, err
	}

	// 3. Mapear a modelo de dominio
	reqInvoice, err := u.mapper.MapToInvoiceData(req, issuer)
	if err != nil {
		return nil, nil, err
	}

	// 4. Crear invoice electrónica a nivel de servicio
	invoice, err := u.invoiceService.Create(ctx, reqInvoice, claims.BranchID)
	if err != nil {
		return nil, nil, err
	}

	// 5. Mapear a modelo de hacienda
	mhInvoice, err := response_mapper.ToMHInvoice(invoice)
	if err != nil {
		return nil, nil, err
	}

	// 6. Configurar detalles de respuesta
	options := &response.SuccessOptions{
		Ambient:        config.Server.AmbientCode,
		GenerationCode: mhInvoice.Identificacion.CodigoGeneracion,
		EmissionDate:   utils.TimeNow(),
	}

	// 7. Comenzar la transmisión de la factura
	result, err := u.transmitter.RetryTransmission(ctx, mhInvoice, token, claims.NIT)
	if err != nil {
		return mhInvoice, options, err
	}
	options.ReceptionStamp = result.ReceptionStamp

	if result.Status != ReceivedStatus {
		logs.Warn("Error transmitting invoice", map[string]interface{}{"error": "TransmissionFailed"})
		return mhInvoice, options, dte_errors.NewDTEErrorSimple("TransmissionFailed")
	}

	// 8. Guardar la factura en la base de datos
	err = u.dteService.Create(ctx, mhInvoice, constants.TransmissionNormal, constants.DocumentReceived, result.ReceptionStamp)
	if err != nil {
		return mhInvoice, options, err
	}

	return mhInvoice, options, nil
}
