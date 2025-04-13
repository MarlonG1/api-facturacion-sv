package retention

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/temporal"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

/*
	1. DocType 1 (Fisico):
		- Tipo de documento document.OperationType
		- Numero de documento (Numero correlativo tradicional) document.NewDocumentNumber(string, documentType)
		- Descripcion string
		- Monto gravado financial.NewAmount
		- Fecha de Emision temporal.NewEmissionDate
		- Tipo de DTE document.NewDTEType
		- Codigo de retencion de MH document.NewRetentionCode

     2. DocType 2 (Electronico):
		- Descripcion string
		- Tipo de documento document.OperationType
		- Tipo de DTE document.NewDTEType
		- Numero de documento (Codigo de generacion UUID) document.NewDocumentNumber(string, documentType)
		- Codigo de retencion de MH document.NewRetentionCode

	Validaciones:
		- Los campos Extension y Appendixes son opcionales en ambos casos, de estar presentes, deben ser validos
		- Si todos los items son document_type 1, entonces la seccion de receptor es obligatoria
		- Si todos los items son document_type 2, entonces la seccion de receptor es excluida
		- Si son mixtos, entonces la seccion receptor es excluida
*/

func MapRetentionItemList(req []structs.RetentionItem) ([]retention_models.RetentionItem, error) {
	if req == nil {
		return nil, dte_errors.NewValidationError("RequiredField", "Items")
	}

	var retentionItems []retention_models.RetentionItem
	for i, item := range req {
		retentionItem, err := mapRetentionItem(&item, i+1)
		if err != nil {
			return nil, err
		}
		retentionItems = append(retentionItems, *retentionItem)
	}

	return retentionItems, nil
}

func mapRetentionItem(req *structs.RetentionItem, i int) (*retention_models.RetentionItem, error) {
	if req == nil {
		return nil, dte_errors.NewValidationError("RequiredField", "RetentionItem")
	}

	documentType, err := document.NewOperationType(req.DocumentType)
	if err != nil {
		return nil, err
	}

	// Formato fisico tradicional
	if documentType.GetValue() == 1 {

		err := validateRetentionPhysicalFields(req)
		if err != nil {
			return nil, err
		}

		documentNumber, err := document.NewDocumentNumber(req.DocumentNumber, req.DocumentType)
		if err != nil {
			return nil, err
		}

		taxedAmount, err := financial.NewAmount(*req.TaxedAmount)
		if err != nil {
			return nil, err
		}

		ivaAmount, err := financial.NewAmount(*req.IvaAmount)
		if err != nil {
			return nil, err
		}

		emissionDate, err := temporal.NewEmissionDateFromString(*req.EmissionDate)
		if err != nil {
			return nil, err
		}

		dteType, err := document.NewDTEType(*req.DTEType)
		if err != nil {
			return nil, err
		}

		retentionCode, err := document.NewRetentionCode(req.RetentionCode)
		if err != nil {
			return nil, err
		}

		return &retention_models.RetentionItem{
			Number:          *item.NewValidatedItemNumber(i),
			DocumentType:    *documentType,
			DocumentNumber:  documentNumber,
			Description:     req.Description,
			RetentionAmount: *taxedAmount,
			RetentionIVA:    *ivaAmount,
			EmissionDate:    *emissionDate,
			DTEType:         *dteType,
			ReceptionCodeMH: *retentionCode,
		}, nil
	}

	documentNumber, err := document.NewDocumentNumber(req.DocumentNumber, req.DocumentType)
	if err != nil {
		return nil, err
	}

	retentionCode, err := document.NewRetentionCode(req.RetentionCode)
	if err != nil {
		return nil, err
	}

	if req.Description == "" {
		return nil, dte_errors.NewValidationError("RequiredField", "Description")
	}

	return &retention_models.RetentionItem{
		Number:          *item.NewValidatedItemNumber(i),
		DocumentType:    *documentType,
		DocumentNumber:  documentNumber,
		Description:     req.Description,
		ReceptionCodeMH: *retentionCode,
	}, nil
}

func validateRetentionPhysicalFields(req *structs.RetentionItem) error {
	if req.TaxedAmount == nil {
		return dte_errors.NewValidationError("RequiredField", "TaxedAmount")
	}

	if req.IvaAmount == nil {
		return dte_errors.NewValidationError("RequiredField", "IVAAmount")
	}

	if req.EmissionDate == nil {
		return dte_errors.NewValidationError("RequiredField", "EmissionDate")
	}

	if req.DTEType == nil {
		return dte_errors.NewValidationError("RequiredField", "DTEType")
	}

	return nil
}
