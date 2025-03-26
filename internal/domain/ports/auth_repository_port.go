package ports

import (
	"context"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
)

// AuthRepositoryPort define el comportamiento que debe implementar un repositorio de autenticación
type AuthRepositoryPort interface {
	// GetAuthTypeByApiKey obtiene el tipo de autenticación de un usuario
	GetAuthTypeByApiKey(context.Context, string) (string, error)
	// GetAuthTypeByNIT obtiene el tipo de autenticación de un usuario
	GetAuthTypeByNIT(context.Context, string) (string, error)
	// GetByNIT obtiene un usuario por su NIT
	GetByNIT(context.Context, string) (*user.User, error)
	// GetIssuerInfoByBranchID obtiene la información del emisor
	GetIssuerInfoByBranchID(context.Context, uint) (*dte.IssuerDTE, error)
	// GetByBranchID obtiene un usuario por su API key de sucursal
	GetByBranchID(context.Context, uint) (*user.User, error)
	// GetBranchByBranchID obtiene la sucursal por su API key
	GetBranchByBranchID(context.Context, uint) (*user.BranchOffice, error)
	// GetBranchByBranchApiKey obtiene la sucursal por su API key
	GetBranchByBranchApiKey(context.Context, string) (*user.BranchOffice, error)
	// GetByBranchApiKey obtiene un usuario por su API key de sucursal
	GetByBranchApiKey(context.Context, string) (*user.User, error)
	// Create crea un usuario con sus sucursales
	Create(context.Context, *user.User) error
	// Update actualiza un usuario
	Update(context.Context, *user.User) error
	// UpdateBranchOffices actualiza las sucursales de un usuario
	UpdateBranchOffices(context.Context, uint, []user.BranchOffice) error
	// DeleteBranchOffice elimina una sucursal de un usuario
	DeleteBranchOffice(context.Context, uint, uint) error
	// GetMatrixBranch obtiene la sucursal registrada como casa matriz
	GetMatrixBranch(context.Context, uint) (*user.BranchOffice, error)
}

// AuthStrategy define el comportamiento que debe implementar cada estrategia de autenticación
type AuthStrategy interface {
	// GetAuthType retorna el tipo de autenticación que implementa esta estrategia
	GetAuthType() string
	// Authenticate valida las credenciales y retorna los claims para el token
	Authenticate(ctx context.Context, credentials *models.AuthCredentials) (*models.AuthClaims, error)
	// ValidateCredentials valida el formato de las credenciales para esta estrategia
	ValidateCredentials(credentials *models.AuthCredentials) error
	// GetHaciendaCredentials obtiene las credenciales de hacienda segun el tipo de autenticación
	GetHaciendaCredentials(token string) (*models.HaciendaCredentials, error)
	// GetTokenLifetime obtiene la duración de vida de un token
	GetTokenLifetime(credentials *models.AuthCredentials) (time.Duration, error)
}

// AuthManager define el comportamiento de un servicio de autenticación
type AuthManager interface {
	// Login maneja el proceso de autenticación
	Login(ctx context.Context, credentials *models.AuthCredentials) (string, error)
	// GetByNIT obtiene un usuario por su NIT
	GetByNIT(ctx context.Context, nit string) (*user.User, error)
	// GetBranchByBranchID obtiene la sucursal por su ID
	GetBranchByBranchID(ctx context.Context, branchID uint) (*user.BranchOffice, error)
	// GetIssuer retorna el emisor por API key
	GetIssuer(ctx context.Context, branchID uint) (*dte.IssuerDTE, error)
	// GetHaciendaCredentials obtiene las credenciales de hacienda segun el tipo de autenticación
	GetHaciendaCredentials(ctx context.Context, nit, token string) (*models.HaciendaCredentials, error)
	// Create crea un usuario con sus sucursales
	Create(ctx context.Context, user *user.User) error
}
