package contingency

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"

	haciendaPorts "github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	authModels "github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/models"
	authPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// ContingencyEventService maneja la preparación y envío de eventos de contingencia
type ContingencyEventService struct {
	authManager  auth.AuthManager
	haciendaAuth haciendaPorts.HaciendaAuthManager
	cache        authPorts.CacheManager
	tokenService authPorts.TokenManager
	signer       haciendaPorts.SignerManager
	repo         contingency.ContingencyRepositoryPort
	timeProvider authPorts.TimeProvider
	httpClient   *http.Client
	connection   *drivers.DbConnection
}

// HaciendaContingencyRequest estructura para la petición de contingencia a Hacienda
type HaciendaContingencyRequest struct {
	NIT      string `json:"nit"`
	Document string `json:"documento"`
}

// NewContingencyEventService constructor para ContingencyEventService
func NewContingencyEventService(
	authManager auth.AuthManager,
	haciendaAuth haciendaPorts.HaciendaAuthManager,
	cache authPorts.CacheManager,
	tokenService authPorts.TokenManager,
	signer haciendaPorts.SignerManager,
	repo contingency.ContingencyRepositoryPort,
	timeProvider authPorts.TimeProvider,
	connection *drivers.DbConnection,
) *ContingencyEventService {
	return &ContingencyEventService{
		authManager:  authManager,
		haciendaAuth: haciendaAuth,
		cache:        cache,
		tokenService: tokenService,
		signer:       signer,
		repo:         repo,
		timeProvider: timeProvider,
		connection:   connection,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

// PrepareAndSendContingencyEvent prepara y envía un evento de contingencia
func (s *ContingencyEventService) PrepareAndSendContingencyEvent(ctx context.Context, docs []dte.ContingencyDocument) error {
	sqlDb, err := s.connection.Db.DB()
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "PrepareAndSendContingencyEvent", "failed to get sql db", err)
	}
	sqlDb.Ping()

	client, err := s.authManager.GetIssuer(ctx, docs[0].BranchID)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "PrepareAndSendContingencyEvent", "failed to get issuer info", err)
	}

	reason, err := s.prepareContingencyReason(docs[0])
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "PrepareAndSendContingencyEvent", "failed to prepare contingency reason", err)
	}

	event := &models.ContingencyEvent{
		Identification: models.ContingencyIdentification{
			Version:          3,
			Ambient:          config.Server.AmbientCode,
			GenerationCode:   strings.ToUpper(uuid.New().String()),
			TransmissionDate: utils.TimeNow().Format("2006-01-02"),
			TransmissionTime: utils.TimeNow().Format("15:04:05"),
		},
		Issuer: models.ContingencyIssuer{
			NIT:                  client.NIT,
			Name:                 client.BusinessName,
			ResponsibleName:      client.BusinessName,
			ResponsibleDocType:   constants.NIT,
			ResponsibleDocNumber: client.NIT,
			EstablishmentType:    client.EstablishmentType,
			Phone:                *client.Phone,
			Email:                *client.Email,
			EstablishmentCodeMH:  client.EstablishmentCodeMH,
			POSCode:              client.POSCode,
		},
		DTEDetails: s.prepareDTEDetails(docs),
		Reason:     reason,
	}

	return s.sendContingencyEvent(ctx, event, docs[0].BranchID)
}

// prepareDTEDetails prepara los detalles de los documentos para el evento de contingencia
func (s *ContingencyEventService) prepareDTEDetails(docs []dte.ContingencyDocument) []models.DTEDetail {
	details := make([]models.DTEDetail, len(docs))
	for i, doc := range docs {
		details[i] = models.DTEDetail{
			ItemNumber:     i + 1,
			GenerationCode: doc.Document.ID,
			DocumentType:   doc.Document.DTEType,
		}
	}
	return details
}

// prepareContingencyReason prepara la razón de contingencia
func (s *ContingencyEventService) prepareContingencyReason(doc dte.ContingencyDocument) (models.ContingencyReason, error) {
	now := s.timeProvider.Now()

	startTime, err := s.repo.GetFirstContingencyTimestamp(context.Background(), doc.BranchID)
	if err != nil {
		return models.ContingencyReason{}, shared_error.NewGeneralServiceError("ContingencyEventService", "prepareContingencyReason", "failed to get first contingency timestamp", err)
	}

	return models.ContingencyReason{
		StartDate:         startTime.Format("2006-01-02"),
		EndDate:           now.Format("2006-01-02"),
		StartTime:         startTime.Format("15:04:05"),
		EndTime:           now.Add(1 * time.Hour).Format("15:04:05"),
		ContingencyType:   doc.ContingencyType,
		ContingencyReason: doc.Reason,
	}, nil
}

// sendContingencyEvent envía el evento de contingencia a Hacienda
func (s *ContingencyEventService) sendContingencyEvent(ctx context.Context, event *models.ContingencyEvent, branchID uint) error {
	// Obtener el client y generar token del sistema
	client, err := s.authManager.GetByNIT(ctx, event.Issuer.NIT)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to get client", err)
	}
	token, err := s.generateMatchingToken(client, branchID)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to generate matching token", err)
	}

	// Obtener credenciales
	encryptedCreds, err := s.cache.GetCredentials(token)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to get hacienda credentials", err)
	}

	// Obtener token de Hacienda
	haciendaToken, err := s.haciendaAuth.GetOrCreateHaciendaTokenWithCreds(
		ctx,
		token,
		*encryptedCreds,
	)

	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to get hacienda token", err)
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to marshal contingency event", err)
	}

	signedDoc, err := s.signer.SignDTE(ctx, jsonData, client.NIT)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to sign contingency event", err)
	}

	reqBody := &HaciendaContingencyRequest{
		NIT:      client.NIT,
		Document: signedDoc,
	}

	jsonData, err = json.Marshal(reqBody)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to marshal contingency request", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.MHPaths.ContingencyURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to create request", err)
	}

	req.Header.Set("Authorization", haciendaToken)
	req.Header.Set("Content-Type", "application/json")

	logs.Info("Sending contingency event request", map[string]interface{}{
		"url":          config.MHPaths.ContingencyURL,
		"method":       "POST",
		"content-type": req.Header.Get("Content-Type"),
	})

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to send request", err)
	}
	defer resp.Body.Close()

	// Manejo de respuesta
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		if err == io.EOF {
			logs.Warn("Contingency event response is empty")
			responseBody = nil
		} else {
			return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to decode response body", err)
		}
	}

	logs.Info("Contingency event response", map[string]interface{}{
		"statusCode": resp.StatusCode,
		"body":       responseBody,
	})

	if resp.StatusCode != http.StatusOK || strings.Contains(responseBody["mensaje"].(string), "no superadas") {
		return shared_error.NewGeneralServiceError("ContingencyEventService", "sendContingencyEvent", "failed to send contingency event", nil)
	}

	return nil
}

// generateMatchingToken genera un token para el cliente
func (s *ContingencyEventService) generateMatchingToken(client *user.User, branchID uint) (string, error) {
	key := fmt.Sprintf("token:timestamps:%d", client.ID)
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
		ClientID: client.ID,
		AuthType: client.AuthType,
		BranchID: branchID,
		NIT:      client.NIT,
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
