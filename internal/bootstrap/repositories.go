package bootstrap

import (
	contiPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/ports"
	dtePorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/repositories"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	db *gorm.DB

	authRepo             ports.AuthRepositoryPort
	sequentialNumberRepo appPorts.SequentialNumberRepositoryPort
	dteRepo              dtePorts.DTERepositoryPort
	contingencyRepo      contiPorts.ContingencyRepositoryPort
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
	c.contingencyRepo = repositories.NewContingencyRepository(c.db)
}

func (c *RepositoryContainer) ContingencyRepo() contiPorts.ContingencyRepositoryPort {
	return c.contingencyRepo
}

func (c *RepositoryContainer) DTERepo() dtePorts.DTERepositoryPort {
	return c.dteRepo
}

func (c *RepositoryContainer) AuthRepo() ports.AuthRepositoryPort {
	return c.authRepo
}

func (c *RepositoryContainer) SequentialNumberRepo() appPorts.SequentialNumberRepositoryPort {
	return c.sequentialNumberRepo
}
