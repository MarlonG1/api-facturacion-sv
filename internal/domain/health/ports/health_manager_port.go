package ports

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
)

// ComponentChecker es una interfaz que define los métodos para verificar el estado de un componente
type ComponentChecker interface {
	Check() models.Health // Check verifica el estado de un componente y devuelve un modelo de models.Health
	Name() string         // Name devuelve el nombre del componente
}

// HealthManager es una interfaz que define los métodos para verificar el estado de todos los componentes
type HealthManager interface {
	CheckHealth() (*models.HealthStatus, error) // CheckHealth verifica el estado de todos los componentes y devuelve un modelo de models.HealthStatus
}
