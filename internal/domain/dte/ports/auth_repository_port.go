package ports

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user/models"

type AuthRepositoryPort interface {
	// GetAuthType obtiene el tipo de autenticación de un usuario
	GetAuthType(uint) (string, error)
	// GetByNIT obtiene un usuario por su NIT
	GetByNIT(string) (*models.User, error)
	// GetUserByBranchOfficeApiKey obtiene un usuario por su API key de sucursal
	GetUserByBranchOfficeApiKey(string) (*models.User, error)
	// Create crea un usuario con sus sucursales
	Create(*models.User, []models.BranchOffice) error
	// Update actualiza un usuario
	Update(*models.User) error
	// UpdateBranchOffices actualiza las sucursales de un usuario
	UpdateBranchOffices(uint, []models.BranchOffice) error
	// DeleteBranchOffice elimina una sucursal de un usuario
	DeleteBranchOffice(uint, uint) error
	// ActivateOrDeactivate activa o desactiva un usuario
	ActivateOrDeactivate(uint, bool) error
}
