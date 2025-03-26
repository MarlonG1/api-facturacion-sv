package service

import (
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
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

func (m *DTEManager) Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error {
	// 1. Establecer el sello de recepción en el apéndice del DTE
	if transmission != constants.TransmissionContingency {
		if err := m.setReceptionStampIntoAppendix(document, receptionStamp); err != nil {
			return shared_error.NewGeneralServiceError("DTEManager", "CreateDTE", "failed to set reception stamp into appendix", err)
		}
	}

	// 2. Crear el DTE en la base de datos
	if err := m.repo.Create(ctx, document, transmission, status, receptionStamp); err != nil {
		return shared_error.NewGeneralServiceError("DTEManager", "CreateDTE", "failed to create DTE", err)
	}

	return nil
}

func (m *DTEManager) UpdateDTE(ctx context.Context, branchID uint, document dte.DTEDetails) error {
	// 1. Actualizar el DTE en la base de datos
	if err := m.repo.Update(ctx, branchID, document); err != nil {
		return shared_error.NewGeneralServiceError("DTEManager", "UpdateDTE", "failed to update DTE", err)
	}

	return nil
}

func (m *DTEManager) VerifyStatus(ctx context.Context, branchID uint, id string) (string, error) {
	// 1. Verificar el estado del DTE en la base de datos
	status, err := m.repo.VerifyStatus(ctx, branchID, id)
	if err != nil {
		return "", shared_error.NewGeneralServiceError("DTEManager", "VerifyStatus", "failed to verify DTE status", err)
	}

	return status, nil
}

func (m *DTEManager) GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error) {
	// 1. Obtener el DTE por su código de generación
	dteDocument, err := m.repo.GetByGenerationCode(ctx, branchID, generationCode)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("DTEManager", "GetByGenerationCode", "failed to get DTE by generation code", err)
	}

	return dteDocument, nil
}

func (m *DTEManager) GetByGenerationCodeConsult(ctx context.Context, branchID uint, generationCode string) (*dte.DTEResponse, error) {
	// 1. Obtener el DTE por su código de generación
	dteDocument, err := m.repo.GetByGenerationCode(ctx, branchID, generationCode)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("DTEManager", "GetByGenerationCodeConsult", "failed to get DTE by generation code", err)
	}

	// 2. Deserializar JSON data
	var jsonData map[string]interface{}
	if err = json.Unmarshal([]byte(dteDocument.Details.JSONData), &jsonData); err != nil {
		return nil, err
	}

	// 3. Retornar la respuesta
	return &dte.DTEResponse{
		GenerationCode: dteDocument.Details.ID,
		ControlNumber:  dteDocument.Details.ControlNumber,
		Status:         dteDocument.Details.Status,
		ReceptionStamp: dteDocument.Details.ReceptionStamp,
		JSONData:       jsonData,
		CreatedAt:      dteDocument.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      dteDocument.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *DTEManager) GetAllDTEs(ctx context.Context, filters *dte.DTEFilters) (*dte.DTEListResponse, error) {
	// 1. Obtener las estadísticas resumidas de los DTEs
	summaryStats, err := m.repo.GetSummaryStats(ctx, filters)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("DTEManager", "GetAll", "failed to get summary stats", err)
	}

	// 2. Crear la respuesta base
	response := &dte.DTEListResponse{
		Summary:   *summaryStats,
		Documents: []dte.DTEModelResponse{},
		Pagination: dte.DTEPaginationResponse{
			Page:       filters.Page,
			PageSize:   filters.PageSize,
			TotalPages: calculateTotalPages(summaryStats.Total, filters.PageSize),
		},
	}

	// 3. Si no hay documentos, retornar la respuesta vacía
	if summaryStats.Total == 0 {
		return response, nil
	}

	// 4. Obtener los documentos paginados
	documents, err := m.repo.GetPagedDocuments(ctx, filters)
	if err != nil {
		return nil, shared_error.NewGeneralServiceError("DTEManager", "GetAll", "failed to get paged documents", err)
	}
	response.Documents = documents

	return response, nil
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

func calculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		pages++
	}

	return pages
}
