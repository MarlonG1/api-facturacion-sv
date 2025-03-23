package bootstrap

import (
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	ports2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/transmission/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/repositories"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	db *gorm.DB

	authRepo             ports.AuthRepositoryPort
	sequentialNumberRepo appPorts.SequentialNumberRepositoryPort
	dteRepo              ports2.DTERepositoryPort
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		db: db,
	}
}

func (c *RepositoryContainer) Initialize() {
	c.authRepo = repositories.NewAuthRepository(c.db)
	c.sequentialNumberRepo = repositories.NewControlNumberRepository(c.db)
	c.dteRepo = repositories.NewDTERepository(c.db)
}

func (c *RepositoryContainer) DTERepo() ports2.DTERepositoryPort {
	return c.dteRepo
}

func (c *RepositoryContainer) AuthRepo() ports.AuthRepositoryPort {
	return c.authRepo
}

func (c *RepositoryContainer) SequentialNumberRepo() appPorts.SequentialNumberRepositoryPort {
	return c.sequentialNumberRepo
}
