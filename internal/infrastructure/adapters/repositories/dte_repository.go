package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/database/db_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
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

func (D DTERepository) Create(ctx context.Context, document interface{}, receptionStamp *string) error {
	// 1. Extraer los claims del contexto
	claims := ctx.Value("claims").(*models.AuthClaims)
	var dteResponse utils.AuxiliarIdentificationExtractor

	// 2. Extraer los datos básicos para el modelo DTEDocument
	jsonData, err := json.Marshal(document)
	if err != nil {
		logs.Error("Failed to marshal DTE", map[string]interface{}{
			"error": err,
			"type":  fmt.Sprintf("%T", document),
		})
		return shared_error.NewGeneralServiceError(
			"BaseDTERepository",
			"Create",
			"failed to marshal DTE",
			err,
		)
	}
	if err := json.Unmarshal(jsonData, &dteResponse); err != nil {
		return shared_error.NewGeneralServiceError(
			"BaseDTERepository",
			"Create",
			"failed to extract DTE fields",
			err,
		)
	}

	// 3. Crear un modelo DTEDocument
	dteDocument := &db_models.DTEDocument{
		BranchID: claims.BranchID,
		Document: &db_models.DTEDetails{
			ID:             dteResponse.Identification.GenerationCode,
			Transmission:   constants.TransmissionNormal,
			Status:         constants.DocumentReceived,
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
