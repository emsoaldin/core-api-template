package roles

import (
	"api/routers"

	"github.com/go-flow/flow/v2"
)

// Module -
type Module struct {
}

// NewModule creates new Roles Module instance
func NewModule() *Module {
	return &Module{}
}

func (m *Module) ProvideImports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewRolesRepository),
	}
}

func (m *Module) ProvideExports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewRolesService),
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
