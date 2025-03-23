package dte

import (
	"context"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	requestDTO "github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type InvoiceUseCase struct {
	authService    ports.AuthManager
	invoiceService interfaces.InvoiceManager
	transmitter    appPorts.BaseTransmitter
	mapper         *request_mapper.InvoiceMapper
}

func NewInvoiceUseCase(authService ports.AuthManager, invoiceService interfaces.InvoiceManager, transmitter appPorts.BaseTransmitter) *InvoiceUseCase {
	return &InvoiceUseCase{
		authService:    authService,
		invoiceService: invoiceService,
		transmitter:    transmitter,
		mapper:         request_mapper.NewInvoiceMapper(),
	}
}

func (u *InvoiceUseCase) Create(ctx context.Context, req *requestDTO.CreateInvoiceRequest) (*structs.InvoiceDTEResponse, *string, error) {
	// 1. Obtener los claims y el token del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	//token := ctx.Value("token").(string)

	logs.Debug("Claims", map[string]interface{}{
		"ID":       claims.ClientID,
		"BranchID": claims.BranchID,
		"NIT":      claims.NIT,
		"AuthType": claims.AuthType,
	})

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

	// 6. Comenzar la transmisión de la factura
	//result, err := u.transmitter.RetryTransmission(ctx, mhInvoice, token, claims.NIT)
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//if result.Status != ReceivedStatus {
	//	return mhInvoice, result.ReceptionStamp, dte_errors.NewDTEErrorSimple("TransmissionFailed")
	//}

	return mhInvoice, new(string), nil
}
