package strategies

import (
	"context"
	"crypto/subtle"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type StandardAuthStrategy struct {
	authRepo     auth.AuthRepositoryPort
	cacheService ports.CacheManager
}

// NewStandardAuthStrategy crea una instancia de StandardAuthStrategy. Recibe un repositorio de clientes.
func NewStandardAuthStrategy(repo auth.AuthRepositoryPort, cacheService ports.CacheManager) *StandardAuthStrategy {
	return &StandardAuthStrategy{
		cacheService: cacheService,
		authRepo:     repo,
	}
}

// GetAuthType devuelve el tipo de autenticación.
func (s *StandardAuthStrategy) GetAuthType() string {
	return constants.StandardAuthType
}

// ValidateCredentials valida las credenciales de autenticación. Devuelve un error si las credenciales son inválidas.
func (s *StandardAuthStrategy) ValidateCredentials(credentials *models.AuthCredentials) error {
	if credentials.APIKey == "" {
		logs.Error("API key is required", map[string]interface{}{
			"credentials": credentials,
		})
		return shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"ValidateCredentials",
			"RequiredField",
			"api key",
		)
	}
	if credentials.APISecret == "" {
		logs.Error("API secret is required", map[string]interface{}{
			"credentials": credentials,
		})
		return shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"ValidateCredentials",
			"RequiredField",
			"api secret",
		)
	}
	return nil
}

// Authenticate autentica un cliente. Devuelve los claims del cliente autenticado.
func (s *StandardAuthStrategy) Authenticate(ctx context.Context, credentials *models.AuthCredentials) (*models.AuthClaims, error) {
	// 1. Obtener sucursal por API key
	branch, err := s.authRepo.GetBranchByBranchApiKey(ctx, credentials.APIKey)
	if err != nil {
		logs.Error("Invalid credentials", map[string]interface{}{
			"apiKey": credentials.APIKey,
			"error":  err.Error(),
		})
		return nil, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"Authenticate",
			"NotFound",
		)
	}

	// 2. Verificar credenciales
	if subtle.ConstantTimeCompare([]byte(credentials.APISecret), []byte(branch.APISecret)) != 1 {
		logs.Error("Invalid credentials", map[string]interface{}{
			"apiKey": credentials.APIKey,
		})
		return nil, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"Authenticate",
			"InvalidCredentials",
		)
	}

	// 3. Obtener usuario por API key
	user, err := s.authRepo.GetByBranchApiKey(ctx, credentials.APIKey)
	if err != nil {
		logs.Error("Invalid credentials", map[string]interface{}{
			"apiKey": credentials.APIKey,
			"error":  err.Error(),
		})
		return nil, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"Authenticate",
			"InvalidCredentials",
		)
	}

	// 4. Verificar estado de cuenta de usuario
	if !user.Status {
		logs.Error("Client account is not active", map[string]interface{}{
			"clientID": user.ID,
		})
		return nil, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"Authenticate",
			"UserNotActive",
		)
	}

	// 5. Crear claims
	claims := &models.AuthClaims{
		ClientID: user.ID,
		BranchID: branch.ID,
		AuthType: user.AuthType,
		NIT:      user.NIT,
	}

	logs.Info("Client authenticated successfully", map[string]interface{}{
		"clientID": claims.ClientID,
	})

	return claims, nil
}

func (s *StandardAuthStrategy) GetTokenLifetime(credentials *models.AuthCredentials) (time.Duration, error) {
	// 1. Obtener informacion del usuario
	user, err := s.authRepo.GetByBranchApiKey(context.Background(), credentials.APIKey)
	if err != nil {
		logs.Error("Failed to get user information", map[string]interface{}{
			"apiKey": credentials.APIKey,
			"error":  err.Error(),
		})
		return 0, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"GetTokenLifetime",
			"ServerError",
		)
	}

	return time.Duration(user.TokenLifetime) * 24 * time.Hour, nil
}

// GetHaciendaCredentials obtiene las credenciales de Hacienda. Devuelve las credenciales de Hacienda.
func (s *StandardAuthStrategy) GetHaciendaCredentials(token string) (*models.HaciendaCredentials, error) {
	creds, err := s.cacheService.GetCredentials(token)
	if err != nil {
		logs.Error("Failed to get Hacienda credentials", map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return nil, shared_error.NewFormattedGeneralServiceError(
			"StandardAuth",
			"GetHaciendaCredentials",
			"FailedToGetCredentials",
		)
	}

	logs.Info("Hacienda credentials retrieved successfully", map[string]interface{}{
		"token": token,
	})

	return creds, nil
}
