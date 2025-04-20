package repositories

import (
	"context"
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"gorm.io/gorm"
)

type FailedSequenceNumberRepository struct {
	db *gorm.DB
}

// NewFailedSequenceNumberRepository crea una nueva instancia de FailedSequenceNumberRepository
func NewFailedSequenceNumberRepository(db *gorm.DB) ports.FailedSequenceNumberRepositoryPort {
	return &FailedSequenceNumberRepository{db: db}
}

// RegisterFailedSequence registra un número de secuencia fallido con detalles
func (r *FailedSequenceNumberRepository) RegisterFailedSequence(
	ctx context.Context,
	branchID uint,
	dteType string,
	sequenceNumber uint,
	year uint,
	failureReason string,
	responseCode string,
	originalRequestData interface{},
	mhResponse string,
) error {
	// Convertir los datos de la solicitud original a JSON
	requestDataJSON, err := json.Marshal(originalRequestData)
	if err != nil {
		logs.Error("Failed to marshal original request data", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Crear la entrada de FailedSequenceNumber
	failedSeq := &db_models.FailedSequenceNumber{
		BranchID:            branchID,
		DTEType:             dteType,
		SequenceNumber:      sequenceNumber,
		Year:                year,
		FailureReason:       failureReason,
		ResponseCode:        responseCode,
		OriginalRequestData: string(requestDataJSON),
		MHResponse:          mhResponse,
		CreatedAt:           utils.TimeNow(),
	}

	result := r.db.WithContext(ctx).Create(failedSeq)
	if result.Error != nil {
		logs.Error("Failed to register failed sequence number", map[string]interface{}{
			"error":          result.Error.Error(),
			"branchID":       branchID,
			"dteType":        dteType,
			"sequenceNumber": sequenceNumber,
		})
		return result.Error
	}

	logs.Info("Failed sequence number registered successfully", map[string]interface{}{
		"id":             failedSeq.ID,
		"branchID":       branchID,
		"dteType":        dteType,
		"sequenceNumber": sequenceNumber,
		"responseCode":   responseCode,
	})

	return nil
}

// GetFailedSequences obtiene una lista de números de secuencia fallidos para una sucursal y tipo de DTE específicos
func (r *FailedSequenceNumberRepository) GetFailedSequences(
	ctx context.Context,
	branchID uint,
	dteType string,
	limit int,
) ([]db_models.FailedSequenceNumber, error) {
	var failedSequences []db_models.FailedSequenceNumber

	result := r.db.WithContext(ctx).
		Where("branch_id = ? AND dte_type = ?", branchID, dteType).
		Order("created_at DESC").
		Limit(limit).
		Find(&failedSequences)

	if result.Error != nil {
		logs.Error("Failed to get failed sequences", map[string]interface{}{
			"error":    result.Error.Error(),
			"branchID": branchID,
			"dteType":  dteType,
		})
		return nil, result.Error
	}

	return failedSequences, nil
}

// GetFailedSequencesByYear obtiene una lista de números de secuencia fallidos para una sucursal, tipo de DTE y año específicos
func (r *FailedSequenceNumberRepository) GetFailedSequencesByYear(
	ctx context.Context,
	branchID uint,
	dteType string,
	year uint,
	limit int,
) ([]db_models.FailedSequenceNumber, error) {
	var failedSequences []db_models.FailedSequenceNumber

	result := r.db.WithContext(ctx).
		Where("branch_id = ? AND dte_type = ? AND year = ?", branchID, dteType, year).
		Order("created_at DESC").
		Limit(limit).
		Find(&failedSequences)

	if result.Error != nil {
		logs.Error("Failed to get failed sequences by year", map[string]interface{}{
			"error":    result.Error.Error(),
			"branchID": branchID,
			"dteType":  dteType,
			"year":     year,
		})
		return nil, result.Error
	}

	return failedSequences, nil
}
