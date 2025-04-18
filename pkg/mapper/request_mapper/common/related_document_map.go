package common

import (
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// MapCommonRequestRelatedDocuments mapea una lista de documentos relacionados a un modelo de documento relacionado -> Origen: Request
func MapCommonRequestRelatedDocuments(relatedDocuments []structs.RelatedDocRequest) ([]models.RelatedDocument, error) {
	result := make([]models.RelatedDocument, len(relatedDocuments))

	for i, relatedDocument := range relatedDocuments {
		relatedDoc, err := MapCommonRequestRelatedDocumentIndex(relatedDocument)
		if err != nil {
			return nil, err
		}
		result[i] = *relatedDoc
	}

	return result, nil
}

// MapCommonRequestRelatedDocumentIndex mapea un documento relacionado a un modelo de documento relacionado -> Origen: Request
func MapCommonRequestRelatedDocumentIndex(doc structs.RelatedDocRequest) (*models.RelatedDocument, error) {

	dteType, err := document.NewDTEType(doc.DocumentType)
	if err != nil {
		return nil, err
	}

	generationType, err := document.NewModelType(doc.GenerationType)
	if err != nil {
		return nil, err
	}

	docNumber, err := document.NewDocumentNumber(doc.DocumentNumber, doc.GenerationType)
	if err != nil {
		return nil, err
	}

	if doc.EmissionDate == "" {
		doc.EmissionDate = utils.TimeNow().Format("2006-01-02")
	}

	timeParse, err := time.Parse("2006-01-02", doc.EmissionDate)
	if err != nil {
		return nil, err
	}

	emissionDate, err := temporal.NewEmissionDate(timeParse)
	if err != nil {
		return nil, err
	}

	return &models.RelatedDocument{
		DocumentType:   *dteType,
		GenerationType: *generationType,
		DocumentNumber: docNumber.GetValue(),
		EmissionDate:   *emissionDate,
	}, nil
}
