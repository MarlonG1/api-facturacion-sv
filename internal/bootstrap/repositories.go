package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/repositories"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	db *gorm.DB

	authRepo ports.AuthRepositoryPort
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		db: db,
	}
}

func (c *RepositoryContainer) Initialize() {
	c.authRepo = repositories.NewAuthRepository(c.db)
}

func (c *RepositoryContainer) AuthRepo() ports.AuthRepositoryPort {
	return c.authRepo
}
