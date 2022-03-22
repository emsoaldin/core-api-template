package auth

import (
	"api/routers"

	"github.com/go-flow/flow/v2"
)

// Module -
type Module struct {
}

// NewModule creates new Auth Module instance
func NewModule() *Module {
	return &Module{}
}

func (m *Module) ProvideImports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewAuthRepository),
	}
}

func (m *Module) ProvideExports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewAuthService),
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
