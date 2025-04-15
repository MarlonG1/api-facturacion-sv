package response_mapper

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/retention"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func ToMHRetention(doc interface{}) *structs.RetentionDTEResponse {

	cast := doc.(*retention_models.RetentionModel)
	dte := &structs.RetentionDTEResponse{
		Identificacion:  common.MapCommonResponseIdentification(cast.Identification),
		Emisor:          retention.MapRetentionResponseIssuer(cast.Issuer),
		Receptor:        common.MapCommonResponseReceiver(cast.Receiver),
		Resumen:         retention.MapRetentionResponseSummary(cast.RetentionSummary),
		CuerpoDocumento: retention.MapRetentionResponseItem(cast.RetentionItems),
		Extension:       common.MapCommonResponseExtension(cast.Extension),
	}

	if cast.Appendix != nil {
		dte.Apendice = common.MapCommonResponseAppendix(cast.Appendix)
	}

	return dte
}
