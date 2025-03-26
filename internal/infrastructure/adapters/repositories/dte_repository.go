package repositories

import (
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"gorm.io/gorm"
)

type DTERepository struct {
	db *gorm.DB
}

func NewDTERepository(db *gorm.DB) ports.DTERepositoryPort {
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
		return "", result.Error
	}

	return status, nil
}

func (D *DTERepository) GetByGenerationCode(ctx context.Context, branchID uint, generationCode string) (*dte.DTEDocument, error) {
	var document db_models.DTEDocument

	// 1. Obtener un documento DTE por código de generación
	result := D.db.WithContext(ctx).
		Preload("Document").
		Where("branch_id = ? AND document_id = ?", branchID, generationCode).
		First(&document)
	if result.Error != nil {
		return nil, result.Error
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
