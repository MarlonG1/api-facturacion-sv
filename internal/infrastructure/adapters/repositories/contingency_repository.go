package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
)

type ContingencyRepository struct {
	db *gorm.DB
}

func NewContingencyRepository(db *gorm.DB) ports.ContingencyRepositoryPort {
	return &ContingencyRepository{db: db}
}

// Create almacena un documento de contingencia en la base de datos
func (r *ContingencyRepository) Create(ctx context.Context, doc *dte.ContingencyDocument) error {
	contingencyDoc := &db_models.ContingencyDocument{
		ID:              uuid.NewString(),
		BranchID:        doc.BranchID,
		DocumentID:      doc.DocumentID,
		ContingencyType: doc.ContingencyType,
		Reason:          doc.Reason,
		CreatedAt:       utils.TimeNow(),
		UpdatedAt:       utils.TimeNow(),
	}

	return r.db.WithContext(ctx).Create(contingencyDoc).Error
}

func (r *ContingencyRepository) GetPending(ctx context.Context, limit int) ([]dte.ContingencyDocument, error) {
	var dbDocs []db_models.ContingencyDocument
	// 1. Obtener los documentos en estado PENDING para procesar (JOIN con dte_details)
	err := r.db.WithContext(ctx).
		Preload("Document").
		Preload("Branch").
		Preload("Branch.User").
		Preload("Branch.Address").
		Joins("JOIN dte_details ON contingency_documents.document_id = dte_details.id").
		Where("dte_details.status = ?", constants.DocumentPending).
		Limit(limit).
		Order("contingency_documents.created_at asc").
		Find(&dbDocs).Error
	if err != nil {
		return nil, err
	}

	// 2. Convertir los documentos a modelos de dominio
	docs := make([]dte.ContingencyDocument, len(dbDocs))
	for i, doc := range dbDocs {
		docs[i] = convertToDomainModel(&doc)
	}

	return docs, nil
}

func (r *ContingencyRepository) UpdateBatch(ctx context.Context, ids []string, observations []string, stamps map[string]string, batchID string, mhBatchID string, status string) error {
	// 1. Iniciar una transacción para asegurar la atomicidad de las operaciones
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Rollback en caso de error
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from panic", r)
			tx.Rollback()
		}
	}()

	for i, id := range ids {
		// 2. Preparar datos básicos de actualización
		contingencyUpdate := map[string]interface{}{
			"batch_id":    batchID,
			"mh_batch_id": mhBatchID,
		}

		// 3. Añadir observaciones si existen para este documento
		if len(observations) > i {
			contingencyUpdate["observations"] = observations[i]
		}

		// 4. Actualizar el documento de contingencia
		if err := tx.Model(&db_models.ContingencyDocument{}).
			Where("id = ?", id).
			Updates(contingencyUpdate).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update document %s: %w", id, err)
		}

		// 5. Obtener el ID del documento asociado a la contingencia
		var contingencyDoc db_models.ContingencyDocument
		if err := tx.Where("id = ?", id).First(&contingencyDoc).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get document %s: %w", id, err)
		}

		// 6. Preparar datos básicos de actualización para el documento asociado
		dteUpdate := map[string]interface{}{
			"status": status,
		}

		// 7. Añadir sello de recepción si existe para este documento
		if stamps != nil {
			if stamp, exists := stamps[id]; exists {
				dteUpdate["reception_stamp"] = stamp

				// 8. Actualizar el apéndice del DTE
				var dteDoc db_models.DTEDetails
				if err := tx.Where("id = ?", contingencyDoc.DocumentID).First(&dteDoc).Error; err != nil {
					logs.Error("Failed to get user", map[string]interface{}{
						"error": err.Error(),
					})
					tx.Rollback()
					return fmt.Errorf("failed to get DTE for appendix update %s: %w", contingencyDoc.DocumentID, err)
				}

				// Actualizamos el apéndice
				updatedJSON, err := utils.SetReceptionStampIntoAppendix(dteDoc.JSONData, &stamp)
				if err != nil {
					logs.Error("Failed to set reception stamp into appendix", map[string]interface{}{
						"error": err.Error(),
					})
					tx.Rollback()
					return fmt.Errorf("failed to update appendix: %w", err)
				}

				// Actualizamos el documento en la misma transacción
				if err := tx.Model(&db_models.DTEDetails{}).
					Where("id = ?", contingencyDoc.DocumentID).
					Update("json_data", updatedJSON).Error; err != nil {
					logs.Error("Failed to update DTE json_data", map[string]interface{}{
						"error": err.Error(),
					})
					tx.Rollback()
					return fmt.Errorf("failed to update DTE json_data %s: %w", contingencyDoc.DocumentID, err)
				}
			}
		}

		// 8. Actualizar el documento asociado
		if err := tx.Model(&db_models.DTEDetails{}).
			Where("id = ?", contingencyDoc.DocumentID).
			Updates(dteUpdate).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update document %s: %w", id, err)
		}
	}

	// 9. Confirmar la transacción si no hay errores
	return tx.Commit().Error
}

func (r *ContingencyRepository) GetFirstContingencyTimestamp(ctx context.Context, branchID uint) (*time.Time, error) {
	var doc db_models.ContingencyDocument

	err := r.db.WithContext(ctx).
		Joins("JOIN dte_details ON contingency_documents.document_id = dte_details.id").
		Where("contingency_documents.branch_id = ? AND dte_details.status = ?", branchID, constants.DocumentPending).
		Order("contingency_documents.created_at ASC").
		Limit(1).
		Select("contingency_documents.created_at").
		First(&doc).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting first contingency timestamp: %w", err)
	}

	return &doc.CreatedAt, nil
}

func convertToDomainModel(doc *db_models.ContingencyDocument) dte.ContingencyDocument {
	return dte.ContingencyDocument{
		ID:              doc.ID,
		BranchID:        doc.BranchID,
		DocumentID:      doc.DocumentID,
		ContingencyType: doc.ContingencyType,
		Reason:          doc.Reason,
		Document: &dte.DTEDetails{
			ID:             doc.Document.ID,
			DTEType:        doc.Document.DTEType,
			ControlNumber:  doc.Document.ControlNumber,
			Transmission:   doc.Document.Transmission,
			Status:         doc.Document.Status,
			ReceptionStamp: doc.Document.ReceptionStamp,
			JSONData:       doc.Document.JSONData,
		},
		Branch: &user.BranchOffice{
			User: &user.User{
				ID:                   doc.Branch.User.ID,
				Status:               doc.Branch.User.Status,
				Email:                doc.Branch.User.Email,
				Phone:                doc.Branch.User.Phone,
				NIT:                  doc.Branch.User.NIT,
				NRC:                  doc.Branch.User.NRC,
				AuthType:             doc.Branch.User.AuthType,
				EconomicActivity:     doc.Branch.User.EconomicActivity,
				EconomicActivityDesc: doc.Branch.User.EconomicActivityDesc,
			},
		},
	}
}
