package service

import (
	"context"
	"errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/service/strategies"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	tokenPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"gorm.io/gorm"
	"strings"
)

type AuthManager struct {
	strategies   map[string]tokenPorts.AuthStrategy
	authRepo     tokenPorts.AuthRepositoryPort
	tokenService tokenPorts.TokenManager
	cacheService tokenPorts.CacheManager
}

func NewAuthService(
	tokenService tokenPorts.TokenManager,
	clientRepository tokenPorts.AuthRepositoryPort,
	cacheService tokenPorts.CacheManager,
) tokenPorts.AuthManager {
	return &AuthManager{
		strategies: map[string]tokenPorts.AuthStrategy{
			constants.StandardAuthType: strategies.NewStandardAuthStrategy(clientRepository, cacheService),
		},
		tokenService: tokenService,
		authRepo:     clientRepository,
		cacheService: cacheService,
	}
}

// Login maneja el proceso de autenticación
func (s *AuthManager) Login(ctx context.Context, credentials *models.AuthCredentials) (string, error) {
	// 0. Verificar existencia de credenciales
	if !credentialsExists(credentials) {
		return "", shared_error.NewGeneralServiceError("AuthManager", "Login", "missing credentials", errors.New("all fields are required"))
	}

	// 1. Obtener tipo de autenticación
	authType, err := s.authRepo.GetAuthTypeByApiKey(ctx, credentials.APIKey)
	if err != nil {
		return "", err
	}

	// 2. Obtener la estrategia apropiada
	strategy, exists := s.strategies[authType]
	if !exists {
		return "", errors.New("unsupported authentication type")
	}

	// 3. Validar formato de credenciales
	if err = strategy.ValidateCredentials(credentials); err != nil {
		return "", err
	}

	// 4. Autenticar usando la estrategia
	claims, err := strategy.Authenticate(ctx, credentials)
	if err != nil {
		return "", err
	}

	// 5. Obtener duración de vida del token
	tokenLifetime, err := strategy.GetTokenLifetime(credentials)
	if err != nil {
		return "", err
	}

	// 5. Generar token JWT
	token, err := s.tokenService.GenerateToken(claims, tokenLifetime)
	if err != nil {
		return "", err
	}

	logs.Debug("TOKEN", map[string]interface{}{"DURATION": tokenLifetime.String()})

	//6. Guardar credenciales en cache
	if err = s.cacheService.SetCredentials(token, credentials.MHCredentials, tokenLifetime); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthManager) GetHaciendaCredentials(ctx context.Context, nit, token string) (*models.HaciendaCredentials, error) {

	// 1. Obtener tipo de autenticación
	authType, err := s.authRepo.GetAuthTypeByNIT(ctx, nit)
	if err != nil {
		logs.Info("Error getting auth type", map[string]interface{}{"error": err.Error()})
		return nil, err
	}
	logs.Info("Auth type retrieved", map[string]interface{}{"authType": authType})

	// 2. Obtener la estrategia apropiada
	strategy, exists := s.strategies[authType]
	if !exists {
		logs.Info("Unsupported authentication type", map[string]interface{}{"authType": authType})
		return nil, errors.New("unsupported authentication type")
	}
	logs.Info("Strategy found", map[string]interface{}{"strategy": strategy.GetAuthType()})

	return strategy.GetHaciendaCredentials(token)
}

// GetIssuer retorna el emisor por su id de sucursal
func (s *AuthManager) GetIssuer(ctx context.Context, branchID uint) (*dte.IssuerDTE, error) {
	return s.authRepo.GetIssuerInfoByBranchID(ctx, branchID)
}

// ValidateToken valida un token existente
func (s *AuthManager) ValidateToken(token string) (*models.AuthClaims, error) {
	return s.tokenService.ValidateToken(token)
}

// RevokeToken revoca un token
func (s *AuthManager) RevokeToken(token string) error {
	return s.tokenService.RevokeToken(token)
}

// credentialsExists verifica que las credenciales tengan todos los campos requeridos
func credentialsExists(credentials *models.AuthCredentials) bool {
	return credentials.APIKey != "" && credentials.APISecret != "" && credentials.MHCredentials != nil && credentials.MHCredentials.Username != "" && credentials.MHCredentials.Password != ""
}

// Create crea un usuario con sus sucursales
func (s *AuthManager) Create(ctx context.Context, user *user.User) error {
	err := s.authRepo.Create(ctx, user)
	if err != nil {
		return handleGormError("create", err)
	}

	return nil
}

func (s *AuthManager) GetByNIT(ctx context.Context, nit string) (*user.User, error) {
	user, err := s.authRepo.GetByNIT(ctx, nit)
	if err != nil {
		return nil, handleGormError("get by nit", err)
	}

	return user, nil
}

func (s *AuthManager) GetBranchByBranchID(ctx context.Context, branchID uint) (*user.BranchOffice, error) {
	branch, err := s.authRepo.GetBranchByBranchID(ctx, branchID)
	if err != nil {
		return nil, handleGormError("get branch by branch id", err)
	}

	return branch, nil
}

func handleGormError(operation string, err error) error {
	if errors.Is(err, gorm.ErrInvalidData) {
		return shared_error.NewGeneralServiceError("AuthService", operation, "invalid data", nil)
	}

	if isDuplicatedEntryErr(err) {
		errMsg := err.Error()
		if strings.Contains(errMsg, "nit") {
			return shared_error.NewGeneralServiceError("AuthService", operation, "nit already exists", nil)
		}

		if strings.Contains(errMsg, "email") {
			return shared_error.NewGeneralServiceError("AuthService", operation, "email already exists", nil)
		}

		if strings.Contains(errMsg, "nrc") {
			return shared_error.NewGeneralServiceError("AuthService", operation, "nrc already exists", nil)
		}
	}

	return err
}

func isDuplicatedEntryErr(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return errors.Is(err, gorm.ErrInvalidData) ||
		strings.Contains(errMsg, "duplicate entry") || // MySQL
		strings.Contains(errMsg, "unique constraint") || // PostgreSQL
		strings.Contains(errMsg, "violates unique") || // PostgreSQL
		strings.Contains(errMsg, "unique key constraint") // SQL Server
}
