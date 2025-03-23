package processors

import (
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type DTEProcessor struct{}

func (p *DTEProcessor) ProcessRequest(signedDoc string, document interface{}) (*models.HaciendaRequest, error) {
	version, dteType, generationCode, sequenceNumber, err := GetDocumentRequestData(document)
	if err != nil {
		logs.Error("Failed to get document request data", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &models.HaciendaRequest{
		Ambient:        env.Server.AmbientCode,
		SendID:         sequenceNumber,
		Version:        version,
		Document:       signedDoc,
		DTEType:        dteType,
		GenerationCode: generationCode,
		URL:            env.MHPaths.ReceptionURL,
	}, nil
}

func (p *DTEProcessor) ProcessResponse(resp *models.HaciendaResponse) (*models.TransmitResult, error) {
	if resp == nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationProcessor", "ProcessResponse", "nil response", nil)
	}

	return &models.TransmitResult{
		Status:         resp.Status,
		ReceptionStamp: &resp.ReceptionStamp,
		ProcessingDate: resp.ProcessingDate,
		MessageCode:    resp.MessageCode,
		MessageDesc:    resp.DescriptionMessage,
		Observations:   resp.Observations,
	}, nil
}
