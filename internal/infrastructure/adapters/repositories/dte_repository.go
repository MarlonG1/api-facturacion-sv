package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type DTERepository struct {
	db *gorm.DB
}

func NewDTERepository(db *gorm.DB) dte_documents.DTERepositoryPort {
	return &DTERepository{
		db: db,
	}
}

func (D *DTERepository) Create(ctx context.Context, document interface{}, transmission, status string, receptionStamp *string) error {
	// 1. Extraer los claims del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	var dteResponse utils.AuxiliarIdentificationExtractor

	// 2. Extraer los datos básicos para el modelo DTEDocument
	jsonData, err := json.Marshal(document)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &dteResponse); err != nil {
		return err
	}

	// 3. Crear un modelo DTEDocument
	dteDocument := &db_models.DTEDocument{
		BranchID:  claims.BranchID,
		CreatedAt: utils.TimeNow(),
		UpdatedAt: utils.TimeNow(),
		Document: &db_models.DTEDetails{
			ID:             dteResponse.Identification.GenerationCode,
			Transmission:   transmission,
			Status:         status,
			DTEType:        dteResponse.Identification.DTEType,
			ControlNumber:  dteResponse.Identification.ControlNumber,
			ReceptionStamp: receptionStamp,
			JSONData:       string(jsonData),
		},
	}

	// 4. Guardar en la base de datos
	result := D.db.WithContext(ctx).Create(dteDocument)
	if result.Error != nil {
		return err
	}

	return nil
}

func (D *DTERepository) GetDTEBalanceControl(ctx context.Context, branchID uint, id string) (*dte.BalanceControl, error) {
	var balanceControl db_models.DTEBalanceControl

	// 1. Obtener el balance de un DTE por su ID
	result := D.db.WithContext(ctx).
		Preload("Transactions").
		Where("branch_id = ? AND original_dte_id = ?", branchID, id).
		First(&balanceControl)
	if result.Error != nil {
		return nil, handleGormErr(result.Error, "GetDTEBalanceControl")
	}

	return &dte.BalanceControl{
		ID:                        balanceControl.ID,
		BranchID:                  balanceControl.BranchID,
		OriginalDTEID:             balanceControl.OriginalDTEID,
		OriginalTaxedAmount:       balanceControl.OriginalTaxedAmount,
		OriginalExemptAmount:      balanceControl.OriginalExemptAmount,
		OriginalNotSubjectAmount:  balanceControl.OriginalTotalNotSubjectAmount,
		RemainingTaxedAmount:      balanceControl.RemainingTaxedAmount,
		RemainingExemptAmount:     balanceControl.RemainingExemptAmount,
		RemainingNotSubjectAmount: balanceControl.RemainingNotSubjectAmount,
	}, nil
}

func (D *DTERepository) GenerateBalanceTransaction(ctx context.Context, branchID uint, originalDTE string, transaction *dte.BalanceTransaction) error {
	var balanceControl db_models.DTEBalanceControl

	// 1. Obtener el balance de un DTE por su ID
	result := D.db.WithContext(ctx).
		Where("branch_id = ? AND original_dte_id = ?", branchID, originalDTE).
		First(&balanceControl)
	if result.Error != nil {
		return handleGormErr(result.Error, "GenerateBalanceTransaction")
	}

	// 2. Crear un nuevo balance de transacción
	dteTransaction := &db_models.DTEBalanceTransaction{
		BalanceControlID:     balanceControl.ID,
		BalanceControl:       &balanceControl,
		AdjustmentDocumentID: transaction.AdjustmentDocumentID,
		TransactionType:      transaction.TransactionType,
		TaxedAmount:          transaction.TaxedAmount,
		ExemptAmount:         transaction.ExemptAmount,
		NotSubjectAmount:     transaction.NotSubjectAmount,
	}

	// 2. Guardar en la base de datos
	result = D.db.WithContext(ctx).Create(dteTransaction)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (D *DTERepository) Update(ctx context.Context, branchID uint, document dte.DTEDetails) error {
	// 1. Actualizar el DTE en la base de datos
	dbModel := &db_models.DTEDetails{
		ID:             document.ID,
		DTEType:        document.DTEType,
		ControlNumber:  document.ControlNumber,
		Transmission:   document.Transmission,
		Status:         document.Status,
		ReceptionStamp: document.ReceptionStamp,
		JSONData:       document.JSONData,
	}

	// 2. Actualizar el DTE en la base de datos
	if err := D.db.WithContext(ctx).
		Model(&db_models.DTEDocument{}).
		Where("document_id = ? AND branch_id = ?", document.ID, branchID).
		Updates(map[string]interface{}{
			"updated_at": utils.TimeNow(),
		}).Error; err != nil {
		return err
	}

	// 3. Actualizar los detalles del DTE
	if err := D.db.WithContext(ctx).
		Model(&db_models.DTEDetails{}).
		Where("id = ?", document.ID).
		Updates(dbModel).Error; err != nil {
		return err
	}

	return nil
}

func (D *DTERepository) VerifyStatus(ctx context.Context, branchID uint, id string) (string, error) {
	var status string

	// 1. Verificar el estado de un DTE en la base de datos
	result := D.db.WithContext(ctx).
		Model(&db_models.DTEDocument{}).
		Joins("JOIN dte_details ON dte_documents.document_id = dte_details.id").
		Select("dte_details.status").
		Where("document_id = ? AND branch_id = ?", id, branchID).
		First(&status)
	if result.Error != nil {
		return "", handleGormErr(result.Error, "VerifyStatus")
	}

	return status, nil
}

// GetTotalCount obtiene el conteo total de documentos que cumplen con los filtros
func (D *DTERepository) GetTotalCount(ctx context.Context, filters *dte.DTEFilters) (int64, error) {
	var totalCount int64

	// 1.Crear la query de consulta en dte_documents junto con dte_details
	query := D.db.WithContext(ctx).
		Table("dte_documents").
		Joins("JOIN dte_details ON dte_documents.document_id = dte_details.id")

	// 2. Aplicar filtros
	loadFilters(query, filters)

	// 3. Ejecutar el conteo
	if err := query.Count(&totalCount).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

// GetSummaryStats obtiene las estadísticas generales para todos los documentos en base a los filtros
func (D *DTERepository) GetSummaryStats(ctx context.Context, filters *dte.DTEFilters) (*dte.ListSummary, error) {
	summary := &dte.ListSummary{}

	// 1. Obtener el conteo total
	totalCount, err := D.GetTotalCount(ctx, filters)
	if err != nil {
		return nil, err
	}
	summary.Total = totalCount

	// 1.1. Si no hay resultados, devolver resumen vacío
	if totalCount == 0 {
		return summary, nil
	}

	// 2. Consulta para agrupar por status y transmission
	query := D.db.WithContext(ctx).
		Table("dte_documents").
		Joins("JOIN dte_details ON dte_documents.document_id = dte_details.id")

	// 3. Aplicar filtros
	loadFilters(query, filters)

	// 4. Estructura para el conteo agrupado
	type StatusTypeCount struct {
		Status       string `gorm:"column:status"`
		Transmission string `gorm:"column:transmission"`
		Count        int64  `gorm:"column:count"`
	}
	var statusCounts []StatusTypeCount

	// 5. Ejecutar la consulta agrupada
	if err := query.Select("dte_details.status, dte_details.transmission, COUNT(*) as count").
		Group("dte_details.status, dte_details.transmission").
		Find(&statusCounts).Error; err != nil {
		return nil, err
	}

	// 6. Mapear conteos a los campos del summary
	for _, sc := range statusCounts {
		// 6.1. Por status
		switch sc.Status {
		case constants.DocumentReceived:
			summary.Received += sc.Count
		case constants.DocumentInvalid:
			summary.Invalid += sc.Count
		case constants.DocumentRejected:
			summary.Rejected += sc.Count
		case constants.DocumentPending:
			summary.Pending += sc.Count
		}

		// 6.2. Por tipo de transmisión
		switch sc.Transmission {
		case constants.TransmissionNormal:
			summary.ByNormal += sc.Count
		case constants.TransmissionContingency:
			summary.ByContingency += sc.Count
		}
	}

	return summary, nil
}

// GetPagedDocuments obtiene los documentos paginados
func (D *DTERepository) GetPagedDocuments(ctx context.Context, filters *dte.DTEFilters) ([]dte.DTEModelResponse, error) {
	// Consulta a dte_documents con preload de dte_details
	query := D.db.WithContext(ctx).
		Table("dte_documents").
		Joins("JOIN dte_details ON dte_documents.document_id = dte_details.id")

	// Aplicar filtros
	loadFilters(query, filters)

	// Estructura para los resultados de la consulta
	type DocumentResult struct {
		ID           string    `gorm:"column:id"`
		JSONData     string    `gorm:"column:json_data"`
		Status       string    `gorm:"column:status"`
		Transmission string    `gorm:"column:transmission"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}

	var documents []DocumentResult

	// Ordenar por fecha de creación (más recientes primero)
	query = query.Order("dte_documents.created_at DESC")

	// Aplicar paginación si es necesario
	if filters.Page > 0 && filters.PageSize > 0 {
		offset := (filters.Page - 1) * filters.PageSize
		query = query.Offset(offset).Limit(filters.PageSize)
	}

	// Seleccionar los campos necesarios
	if err := query.Select("dte_details.id, dte_details.json_data, dte_details.status, dte_details.transmission, dte_documents.created_at").
		Find(&documents).Error; err != nil {
		return nil, err
	}

	// Mapear a DTEModelResponse
	result := make([]dte.DTEModelResponse, 0, len(documents))
	for _, doc := range documents {
		result = append(result, dte.DTEModelResponse{
			Status:           doc.Status,
			TransmissionType: doc.Transmission,
			Document:         json.RawMessage(doc.JSONData),
		})
	}

	return result, nil
}

func (D *DTERepository) GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error) {
	var document db_models.DTEDocument

	// 1. Obtener un documento DTE por código de generación
	result := D.db.WithContext(ctx).
		Preload("Document").
		Where("branch_id = ? AND document_id = ?", branchID, generationCode).
		First(&document)
	if result.Error != nil {
		return nil, handleGormErr(result.Error, "GetByGenerationCode")
	}

	// 3. Retornar el documento DTE
	return &dte.DTEDocument{
		DocumentID: document.Document.ID,
		BranchID:   document.BranchID,
		CreatedAt:  document.CreatedAt,
		UpdatedAt:  document.UpdatedAt,
		Details: &dte.DTEDetails{
			ID:             document.Document.ID,
			DTEType:        document.Document.DTEType,
			ControlNumber:  document.Document.ControlNumber,
			Transmission:   document.Document.Transmission,
			Status:         document.Document.Status,
			ReceptionStamp: document.Document.ReceptionStamp,
			JSONData:       document.Document.JSONData,
		},
	}, nil
}

func loadFilters(query *gorm.DB, filters *dte.DTEFilters) {
	if filters.BranchID != 0 {
		query = query.Where("dte_documents.branch_id = ?", filters.BranchID)
	}

	if filters.DTEType != "" {
		query = query.Where("dte_details.dte_type = ?", filters.DTEType)
	}

	if filters.Status != "" {
		query = query.Where("dte_details.status = ?", filters.Status)
	}

	if filters.StartDate != nil {
		query = query.Where("dte_documents.created_at >= ?", filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("dte_documents.created_at <= ?", filters.EndDate)
	}

	if filters.Transmission != "" {
		query = query.Where("dte_details.transmission = ?", filters.Transmission)
	}
}

func handleGormErr(err error, operation string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared_error.NewFormattedGeneralServiceError("DTERepo", operation, "NotFound")
	}

	return err
}
