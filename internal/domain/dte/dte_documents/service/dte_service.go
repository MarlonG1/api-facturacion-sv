package service

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte_documents/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type DTEManager struct {
	repo ports.DTERepositoryPort
}

func NewDTEManager(repo ports.DTERepositoryPort) interfaces.DTEManager {
	return &DTEManager{
		repo: repo,
	}
}

func (m *DTEManager) Create(ctx context.Context, document interface{}, receptionStamp *string) error {
	// 1. Establecer el sello de recepción en el apéndice del DTE
	if err := m.setReceptionStampIntoAppendix(document, receptionStamp); err != nil {
		return shared_error.NewGeneralServiceError("DTEManager", "CreateDTE", "failed to set reception stamp into appendix", err)
	}

	// 2. Crear el DTE en la base de datos
	if err := m.repo.Create(ctx, document, receptionStamp); err != nil {
		return shared_error.NewGeneralServiceError("DTEManager", "CreateDTE", "failed to create DTE", err)
	}

	return nil
}

func (m *DTEManager) setReceptionStampIntoAppendix(document interface{}, receptionStamp *string) error {
	// 1. Determinar el tipo de DTE
	dteType, err := m.determineDTEType(document)
	if err != nil || dteType == "" {
		return shared_error.NewGeneralServiceError("DTEManager", "setReceptionStampIntoAppendix", "failed to determine DTE type", nil)
	}

	// 2. Deserializar en el modelo correspondiente
	appendix := &structs.DTEApendice{
		Campo:    "Datos del documento",
		Etiqueta: "Sello de recepción",
		Valor:    *receptionStamp,
	}

	// 3. Agregar el sello de recepción al apéndice
	switch dteType {
	case constants.FacturaElectronica:
		document.(*structs.InvoiceDTEResponse).Apendice =
			append(document.(*structs.InvoiceDTEResponse).Apendice, *appendix)
	case constants.CCFElectronico:
		document.(*structs.CCFDTEResponse).Apendice =
			append(document.(*structs.CCFDTEResponse).Apendice, *appendix)
	}

	return nil
}

func (m *DTEManager) determineDTEType(document interface{}) (string, error) {
	dteExtracted, err := utils.ExtractAuxiliarIdentification(document)

	if err != nil {
		return "", shared_error.NewGeneralServiceError("DTEManager", "determineDTEType", "failed to extract DTE identification", err)
	}

	return dteExtracted.Identification.DTEType, nil
}
