package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/retention"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToRetentionMH(doc *retention_models.RetentionModel) (*structs.RetentionDTEResponse, error) {
	dte := &structs.RetentionDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(doc.Identification),
		Emisor:          retention.MapRetentionResponseIssuer(doc.Issuer),
		Receptor:        common.MapCommonResponseReceiver(doc.Receiver),
		Resumen:         retention.MapRetentionResponseSummary(doc.RetentionSummary),
		CuerpoDocumento: retention.MapRetentionResponseItem(doc.RetentionItems),
		Extension:       common.MapCommonResponseExtension(doc.Extension),
	}

	if doc.Appendix != nil {
		dte.Apendice = common.MapCommonResponseAppendix(doc.Appendix)
	}

	return dte, nil
}
