package users

import (
	"api/modules/users/repositories"
	"api/modules/users/routers"
	"api/modules/users/services"

	"github.com/go-flow/flow/v2"
)

// Module -
type Module struct {
}

// NewModule creates new User Module instance
func NewModule() *Module {
	return &Module{}
}

func (m *Module) ProvideImports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(repositories.NewUsersRepository),
	}
}

func (m *Module) ProvideExports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(services.NewUsersService),
	}
}

func (m *Module) ProvideModules() []flow.Provider {
	return []flow.Provider{}
}

func (m *Module) ProvideRouters() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(routers.NewRouter),
	}
}
