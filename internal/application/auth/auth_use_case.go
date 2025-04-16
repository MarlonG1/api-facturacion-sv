package auth

import (
	"context"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type AuthUseCase struct {
	authManager  auth.AuthManager
	cryptManager ports.CryptManager
}

func NewAuthUseCase(authManager auth.AuthManager, cryptManager ports.CryptManager) *AuthUseCase {
	return &AuthUseCase{
		authManager:  authManager,
		cryptManager: cryptManager,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, credentials *models.AuthCredentials) (string, error) {
	// 1. Validar las credenciales obtenidas del request
	if err := credentials.Validate(); err != nil {
		return "", err
	}

	// 2. Autenticar al usuario
	token, err := a.authManager.Login(ctx, credentials)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthUseCase) Register(ctx context.Context, user *user.User) ([]user.ListBranchesResponse, error) {
	// 1. Validar los datos del usuario
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 2.Generar todas las API KEYS y API SECRETS necesarios para el usuario
	keys, secrets, err := a.cryptManager.GenerateBulkAPIKeys(len(user.BranchOffices))
	if err != nil {
		logs.Error("Failed to generate bulk API keys and secrets", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, shared_error.NewFormattedGeneralServiceError("AuthUseCase", "Register", "FailedToCreateUser")
	}

	// 3. Asignar las API KEYS y API SECRETS a las sucursales del usuario
	user.SetBranchesKeysAndSecrets(keys, secrets)

	//4. Crear el usuario en la base de datos
	if err = a.authManager.Create(ctx, user); err != nil {
		logs.Error("Failed to create user", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, shared_error.NewFormattedGeneralServiceWithError("AuthUseCase", "Register", err, "FailedToCreateUser")
	}

	return user.ListBranches(), nil
}
