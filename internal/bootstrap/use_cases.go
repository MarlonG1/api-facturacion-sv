package bootstrap

import "github.com/MarlonG1/api-facturacion-sv/internal/application/auth"

type UseCaseContainer struct {
	services *ServicesContainer

	authUseCase *auth.AuthUseCase
}

func NewUseCaseContainer(services *ServicesContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() {
	c.authUseCase = auth.NewAuthUseCase(c.services.AuthManager(), c.services.CryptManager())
}

func (c *UseCaseContainer) AuthUseCase() *auth.AuthUseCase {
	return c.authUseCase
}
