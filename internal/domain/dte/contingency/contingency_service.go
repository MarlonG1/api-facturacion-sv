package contingency

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/config"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	authModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents"
	batch "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter"
	transmitterModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type ContingencyService struct {
	authManager       auth.AuthManager
	dteManager        dte_documents.DTEManager
	repo              ContingencyRepositoryPort
	haciendaAuth      appPorts.HaciendaAuthManager
	cache             ports.CacheManager
	tokenService      ports.TokenManager
	signer            appPorts.SignerManager
	batchTransmitter  batch.BatchTransmitterPort
	contingencyEvents ContingencyEventSender
	timeProvider      ports.TimeProvider
	config            *transmitterModels.TransmissionConfig
}

func NewContingencyManager(
	authManager auth.AuthManager,
	dteManager dte_documents.DTEManager,
	repo ContingencyRepositoryPort,
	haciendaAuth appPorts.HaciendaAuthManager,
	cache ports.CacheManager,
	tokenService ports.TokenManager,
	signer appPorts.SignerManager,
	batchTransmitter batch.BatchTransmitterPort,
	contingencyEvents ContingencyEventSender,
	timeProvider ports.TimeProvider,
	config *transmitterModels.TransmissionConfig,
) ContingencyManager {
	return &ContingencyService{
		authManager:       authManager,
		dteManager:        dteManager,
		repo:              repo,
		haciendaAuth:      haciendaAuth,
		cache:             cache,
		tokenService:      tokenService,
		signer:            signer,
		batchTransmitter:  batchTransmitter,
		contingencyEvents: contingencyEvents,
		config:            config,
		timeProvider:      timeProvider,
	}
}

// StoreDocumentInContingency almacena un documento en contingencia
func (s *ContingencyService) StoreDocumentInContingency(ctx context.Context, document interface{}, dteType string, contingencyType int8, reason string) error {
	// 1. Extraer los claims del contexto
	claims := ctx.Value("claims").(*authModels.AuthClaims)

	// 2. Extraer información general del documento
	dteInfo, err := utils.ExtractAuxiliarIdentification(document)
	if err != nil {
		logs.Error("Failed to extract general DTE info", map[string]interface{}{
			"error": err.Error(),
			"type":  dteType,
		})
		return shared_error.NewGeneralServiceError("ContingencyService", "StoreDocumentInContingency", "failed to extract general DTE info", err)
	}

	// 3. Almacenar el documento en la base de datos
	err = s.dteManager.Create(ctx, document, constants.TransmissionContingency,
		constants.DocumentPending, nil)
	if err != nil {
		logs.Error("Failed to store DTE", map[string]interface{}{
			"error": err.Error(),
			"type":  dteType,
		})
		return shared_error.NewGeneralServiceError("ContingencyService", "StoreDocumentInContingency", "failed to store DTE", err)
	}

	// 4. Generar el documento de contingencia
	contingencyDoc := &dte.ContingencyDocument{
		DocumentID:      dteInfo.Identification.GenerationCode,
		BranchID:        claims.BranchID,
		ContingencyType: contingencyType,
		Reason:          reason,
	}

	// 5. Almacenar el documento en contingencia
	if err = s.repo.Create(ctx, contingencyDoc); err != nil {
		logs.Error("Failed to store contingency document", map[string]interface{}{
			"error": err.Error(),
			"id":    contingencyDoc.ID,
		})
		return shared_error.NewGeneralServiceError("ContingencyService", "StoreDocumentInContingency", "failed to store contingency document", err)
	}

	logs.Info("Document stored in contingency", map[string]interface{}{
		"id":              contingencyDoc.ID,
		"type":            dteType,
		"contingencyType": contingencyType,
	})

	return nil
}

// RetransmitPendingDocuments retransmite documentos pendientes en contingencia
func (s *ContingencyService) RetransmitPendingDocuments(ctx context.Context) error {
	pendingDocs, err := s.repo.GetPending(ctx, config.Server.MaxBatchSize)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyService", "RetransmitPendingDocuments", "failed to get pending documents", err)
	}

	if len(pendingDocs) == 0 {
		logs.Info("No pending documents found")
		return nil
	}

	// Agrupar por sistema y tipo de DTE
	docsBySystemAndType := s.groupBySystemAndType(pendingDocs)

	for systemNIT, typeGroups := range docsBySystemAndType {
		// Primero enviar el evento de contingencia para todos los documentos del sistema
		if err := s.contingencyEvents.PrepareAndSendContingencyEvent(ctx, pendingDocs); err != nil {
			logs.Error("Failed to send contingency event", map[string]interface{}{
				"error":     err.Error(),
				"systemNIT": systemNIT,
			})
			continue
		}

		// Luego procesar cada grupo de documentos por tipo
		for dteType, docs := range typeGroups {
			if err := s.processSystemDocumentsByType(ctx, systemNIT, dteType, docs); err != nil {
				logs.Error("Failed to process system documents", map[string]interface{}{
					"error":     err.Error(),
					"systemNIT": systemNIT,
					"dteType":   dteType,
				})
				continue
			}
		}
	}

	return nil
}

// processSystemDocumentsByType procesa documentos de un tipo específico para un sistema
func (s *ContingencyService) processSystemDocumentsByType(ctx context.Context, systemNIT string, dteType string, docs []dte.ContingencyDocument) error {
	if len(docs) == 0 {
		logs.Warn("No documents to process")
		return nil
	}
	// 1. Obtener el cliente y generar token
	branchID := docs[0].BranchID
	client, err := s.authManager.GetBranchByBranchID(ctx, branchID)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyService", "processSystemDocumentsByType", "failed to get branch by ID", err)
	}
	token, err := s.generateMatchingToken(client)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyService", "processSystemDocumentsByType", "failed to generate matching token", err)
	}

	// 2. Obtener credenciales
	encryptedCreds, err := s.cache.GetCredentials(token)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyService", "processSystemDocumentsByType", "failed to get credentials", err)
	}

	// 3. Procesar documentos en lotes de máximo 100
	for i := 0; i < len(docs); i += s.config.GetBatchSize() {
		end := i + s.config.GetBatchSize()
		if end > len(docs) {
			end = len(docs)
		}
		batchDocs := docs[i:end]

		// Firmar documentos del lote
		signedDocs := make([]string, 0)
		docsMap := make(map[string]dte.ContingencyDocument)
		docIds := make([]string, 0)

		for _, doc := range batchDocs {
			signedDoc, err := s.signer.SignDTE(ctx, []byte(doc.Document.JSONData), systemNIT)
			if err != nil {
				logs.Error("Failed to sign document", map[string]interface{}{
					"error": err.Error(),
					"nit":   systemNIT,
					"id":    doc.ID,
				})
				continue
			}

			docsMap[doc.Document.ID] = doc
			docIds = append(docIds, doc.ID)
			signedDocs = append(signedDocs, signedDoc)
		}

		if len(signedDocs) == 0 {
			logs.Warn("No documents signed")
			continue
		}

		// Enviar el lote
		batchID := strings.ToUpper(uuid.New().String())

		// Transmitir el lote
		response, haciendaToken, err := s.batchTransmitter.TransmitBatch(ctx, systemNIT, dteType, signedDocs, token, *encryptedCreds)
		if err != nil {
			logs.Error("Failed to transmit batch", map[string]interface{}{
				"error":    err.Error(),
				"batchId":  batchID,
				"dteType":  dteType,
				"docCount": len(signedDocs),
			})
			continue
		}

		// Verificar el estado del lote y procesar resultados
		err = s.batchTransmitter.VerifyContingencyBatchStatus(ctx, batchID, response.BatchCode, haciendaToken, branchID, docsMap)
		if err != nil {
			logs.Error("Failed to verify batch status", map[string]interface{}{
				"error":   err.Error(),
				"batchId": batchID,
				"dteType": dteType,
			})
		}
	}

	return nil
}

// groupBySystemAndType agrupa documentos por sistema y tipo
func (s *ContingencyService) groupBySystemAndType(docs []dte.ContingencyDocument) map[string]map[string][]dte.ContingencyDocument {
	result := make(map[string]map[string][]dte.ContingencyDocument)
	for _, doc := range docs {
		if result[doc.Branch.User.NIT] == nil {
			result[doc.Branch.User.NIT] = make(map[string][]dte.ContingencyDocument)
		}
		result[doc.Branch.User.NIT][doc.Document.DTEType] = append(result[doc.Branch.User.NIT][doc.Document.DTEType], doc)
	}
	return result
}

// generateMatchingToken genera un token para el cliente
func (s *ContingencyService) generateMatchingToken(client *user.BranchOffice) (string, error) {
	key := fmt.Sprintf("token:timestamps:%d", client.User.ID)
	var timestamps struct {
		IssuedAt  int64 `json:"IssuedAt"`
		ExpiresAt int64 `json:"ExpiresAt"`
	}

	jsonTimestamps, err := s.cache.Get(key)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(jsonTimestamps), &timestamps); err != nil {
		return "", err
	}

	claims := &authModels.AuthClaims{
		ClientID: client.User.ID,
		BranchID: client.ID,
		AuthType: client.User.AuthType,
		NIT:      client.User.NIT,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        claims.ClientID,
		"branch_sub": claims.BranchID,
		"auth_type":  claims.AuthType,
		"nit":        claims.NIT,
		"exp":        timestamps.ExpiresAt,
		"iat":        timestamps.IssuedAt,
	})

	return token.SignedString([]byte(s.tokenService.GetSecretKey()))
}
