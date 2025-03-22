package bootstrap

import (
	"gorm.io/gorm"
	"sync"
)

type Container struct {
	db           *gorm.DB
	repositories *RepositoryContainer
	services     *ServicesContainer
	useCases     *UseCaseContainer
	handlers     *HandlerContainer
	middleware   *MiddlewareContainer
	mu           sync.RWMutex
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		db: db,
	}
}

// Initialize inicializa todos los contenedores de la aplicación
func (c *Container) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Inicializar containers en orden de dependencia
	c.repositories = NewRepositoryContainer(c.db)
	c.repositories.Initialize()

	c.services = NewServicesContainer(c.repositories)
	if err := c.services.Initialize(); err != nil {
		return err
	}

	c.useCases = NewUseCaseContainer(c.services)
	c.useCases.Initialize()

	c.middleware = NewMiddlewareContainer(c.services)
	c.middleware.Initialize()

	c.handlers = NewHandlerContainer(c.useCases)
	c.handlers.Initialize()

	return nil
}

func (c *Container) Repositories() *RepositoryContainer {
	return c.repositories
}

func (c *Container) Services() *ServicesContainer {
	return c.services
}

func (c *Container) UseCases() *UseCaseContainer {
	return c.useCases
}

func (c *Container) Handlers() *HandlerContainer {
	return c.handlers
}

func (c *Container) Middleware() *MiddlewareContainer {
	return c.middleware
}
