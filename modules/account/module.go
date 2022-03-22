package account

import (
	"api/modules/account/routers"
	"api/modules/account/services"
	"api/modules/auth"
	"api/modules/roles"
	"api/modules/tokens"
	"api/modules/users"

	"github.com/go-flow/flow/v2"
)

// Module -
type Module struct {
}

// NewModule creates New AccountModule instance
func NewModule() *Module {
	return &Module{}
}

func (m *Module) ProvideImports() []flow.Provider {
	return []flow.Provider{}
}

func (m *Module) ProvideExports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(services.NewAccountService),
	}
}

func (m *Module) ProvideModules() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(roles.NewModule),
		flow.NewProvider(users.NewModule),
		flow.NewProvider(auth.NewModule),
		flow.NewProvider(tokens.NewModule),
	}
}

func (m *Module) ProvideRouters() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(routers.NewPublicRouter),
	}
}
