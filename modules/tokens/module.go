package tokens

import "github.com/go-flow/flow/v2"

// Module -
type Module struct {
}

// NewModule creates new Tokens Module instance
func NewModule() *Module {
	return &Module{}
}

func (m *Module) ProvideImports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewTokensRepository),
	}
}

func (m *Module) ProvideExports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(NewTokensService),
	}
}

func (m *Module) ProvideModules() []flow.Provider {
	return []flow.Provider{}
}

func (m *Module) ProvideRouters() []flow.Provider {
	return []flow.Provider{}
}
