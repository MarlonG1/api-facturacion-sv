package dte_documents

import (
	"context"
	"encoding/json"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type DTEService struct {
	repo DTERepositoryPort
}

func NewDTEService(repo DTERepositoryPort) DTEManager {
	return &DTEService{
		repo: repo,
	}
}

func (m *DTEService) Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error {
	// 1. Establecer el sello de recepción en el apéndice del DTE
	if transmission != constants.TransmissionContingency {
		if err := m.setReceptionStampIntoAppendix(document, receptionStamp); err != nil {
			return shared_error.NewFormattedGeneralServiceWithError("DTEService", "CreateDTE", err, "FailedToSetReceptionStamp")
		}
	}

	// 2. Crear el DTE en la base de datos
	if err := m.repo.Create(ctx, document, transmission, status, receptionStamp); err != nil {
		return shared_error.NewFormattedGeneralServiceWithError("DTEService", "CreateDTE", err, "FailedToCreateDTE")
	}

	return nil
}

func (m *DTEService) GenerateBalanceTransaction(ctx context.Context, branchID uint, transactionType, originalDTE, adjustmentDTE string, document interface{}) error {
	// 1. Extracer los datos del DTE
	extractor, err := utils.ExtractSummaryTotalAmounts(document)
	if err != nil {
		return shared_error.NewFormattedGeneralServiceError("DTEService", "GenerateBalanceTransaction", "FailedToExtractSummaryTotals")
	}

	// 2. Crear la transacción de balance
	transaction := dte.BalanceTransaction{
		AdjustmentDocumentID: adjustmentDTE,
		TransactionType:      transactionType,
		TaxedAmount:          extractor.Summary.TotalTaxed,
		ExemptAmount:         extractor.Summary.TotalExempt,
		NotSubjectAmount:     extractor.Summary.TotalNotSubject,
	}

	// 3. Generar el control de saldo del DTE
	err = m.repo.GenerateBalanceTransaction(ctx, branchID, originalDTE, &transaction)
	if err != nil {
		return shared_error.NewFormattedGeneralServiceWithError("DTEService", "GenerateBalanceTransaction", err, "FailedToGenerateBalanceTransaction")
	}

	return nil
}

func (m *DTEService) GenerateBalanceTransactionWithAmounts(ctx context.Context, branchID uint, transactionType, originalDTE, adjustmentDTE string, taxedSale, exemptSale, notSubjectSale float64) error {
	// 1. Crear la transacción de balance
	transaction := dte.BalanceTransaction{
		AdjustmentDocumentID: adjustmentDTE,
		TransactionType:      transactionType,
		TaxedAmount:          taxedSale,
		ExemptAmount:         exemptSale,
		NotSubjectAmount:     notSubjectSale,
	}

	// 2. Generar el control de saldo del DTE
	err := m.repo.GenerateBalanceTransaction(ctx, branchID, originalDTE, &transaction)
	if err != nil {
		return shared_error.NewFormattedGeneralServiceWithError("DTEService", "GenerateBalanceTransactionWithAmounts", err, "FailedToGenerateBalanceTransaction")
	}

	return nil
}

func (m *DTEService) ValidateForCreditNote(ctx context.Context, branchID uint, originalDTE string, document interface{}) error {

	doc := document.(*credit_note_models.CreditNoteModel)
	if doc == nil {
		return shared_error.NewGeneralServiceError("DTEService", "ValidateForCreditNote", "failed to cast document to CreditNoteModel", nil)
	}

	// 2. Obtener el control de saldo del DTE
	balanceControl, err := m.repo.GetDTEBalanceControl(ctx, branchID, originalDTE)
	if err != nil {
		return shared_error.NewFormattedGeneralServiceWithError("DTEService", "IsValidForCreditNote", err, "FailedToGetBalanceControl")
	}

	// 3. Verificar si el DTE es válido para la Nota de Crédito
	if (balanceControl.RemainingTaxedAmount - doc.Summary.GetTotalTaxed()) < 0 {
		return dte_errors.NewValidationError("InvalidCreditNoteTransaction", "Taxed", originalDTE, doc.Summary.GetTotalTaxed(), balanceControl.RemainingTaxedAmount)
	}
	if (balanceControl.RemainingExemptAmount - doc.Summary.GetTotalExempt()) < 0 {
		return dte_errors.NewValidationError("InvalidCreditNoteTransaction", "Exempt", originalDTE, doc.Summary.GetTotalExempt(), balanceControl.RemainingExemptAmount)
	}
	if (balanceControl.RemainingNotSubjectAmount - doc.Summary.GetTotalNonSubject()) < 0 {
		return dte_errors.NewValidationError("InvalidCreditNoteTransaction", "Not Subject", originalDTE, doc.Summary.GetTotalNonSubject(), balanceControl.RemainingNotSubjectAmount)
	}

	return nil
}

func (m *DTEService) UpdateDTE(ctx context.Context, branchID uint, document dte.DTEDetails) error {
	// 1. Actualizar el DTE en la base de datos
	if err := m.repo.Update(ctx, branchID, document); err != nil {
		return shared_error.NewFormattedGeneralServiceWithError("DTEService", "UpdateDTE", err, "FailedToUpdateDTE")
	}

	return nil
}

func (m *DTEService) VerifyStatus(ctx context.Context, branchID uint, id string) (string, error) {
	// 1. Verificar el estado del DTE en la base de datos
	status, err := m.repo.VerifyStatus(ctx, branchID, id)
	if err != nil {
		return "", shared_error.NewFormattedGeneralServiceWithError("DTEService", "VerifyStatus", err, "FailedToVerifyDTE")
	}

	return status, nil
}

func (m *DTEService) GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error) {
	// 1. Obtener el DTE por su código de generación
	dteDocument, err := m.repo.GetByGenerationCode(ctx, branchID, generationCode)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("DTEService", "GetByGenerationCode", err, "FailedToGetDTE", generationCode)
	}

	return dteDocument, nil
}

func (m *DTEService) GetByGenerationCodeConsult(ctx context.Context, branchID uint, generationCode string) (*dte.DTEResponse, error) {
	// 1. Obtener el DTE por su código de generación
	dteDocument, err := m.repo.GetByGenerationCode(ctx, branchID, generationCode)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceWithError("DTEService", "GetByGenerationCode", err, "FailedToGetDTE", generationCode)
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
		Transmission:   dteDocument.Details.Transmission,
		ReceptionStamp: dteDocument.Details.ReceptionStamp,
		JSONData:       jsonData,
		CreatedAt:      dteDocument.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      dteDocument.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *DTEService) GetAllDTEs(ctx context.Context, filters *dte.DTEFilters) (*dte.DTEListResponse, error) {
	// 1. Obtener las estadísticas resumidas de los DTEs
	summaryStats, err := m.repo.GetSummaryStats(ctx, filters)
	if err != nil {
		return nil, shared_error.NewFormattedGeneralServiceError("DTEService", "GetAll", "FailedToGetSummaryStats")
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
		return nil, shared_error.NewFormattedGeneralServiceError("DTEService", "GetAll", "FailedToGetPagedDoc")
	}
	response.Documents = documents

	return response, nil
}

func (m *DTEService) setReceptionStampIntoAppendix(document interface{}, receptionStamp *string) error {
	// 1. Determinar el tipo de DTE
	dteType, err := m.determineDTEType(document)
	if err != nil || dteType == "" {
		return shared_error.NewGeneralServiceError("DTEService", "setReceptionStampIntoAppendix", "failed to determine DTE type", nil)
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
	case constants.NotaCreditoElectronica:
		document.(*structs.CreditNoteDTEResponse).Apendice =
			append(document.(*structs.CreditNoteDTEResponse).Apendice, *appendix)
	case constants.ComprobanteRetencionElectronico:
		document.(*structs.RetentionDTEResponse).Apendice =
			append(document.(*structs.RetentionDTEResponse).Apendice, *appendix)
	}

	return nil
}

func (m *DTEService) determineDTEType(document interface{}) (string, error) {
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
