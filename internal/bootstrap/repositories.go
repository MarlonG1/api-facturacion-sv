package bootstrap

import (
	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	contiPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/ports"
	dtePorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/dte_documents/ports"
	appPorts "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/adapters/repositories"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	connection *drivers.DbConnection
	db         *gorm.DB

	authRepo             ports.AuthRepositoryPort
	sequentialNumberRepo appPorts.SequentialNumberRepositoryPort
	dteRepo              dtePorts.DTERepositoryPort
	contingencyRepo      contiPorts.ContingencyRepositoryPort
}

func NewRepositoryContainer(connection *drivers.DbConnection) *RepositoryContainer {
	return &RepositoryContainer{
		connection: connection,
		db:         connection.Db,
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
