package routers

import (
	"api/modules/account/actions"

	"github.com/go-flow/flow/v2"
)

// PublicRouter
type PublicRouter struct {
}

func NewPublicRouter() *PublicRouter {
	return &PublicRouter{}
}

func (r *PublicRouter) Path() string {
	return "/account"
}

func (r *PublicRouter) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

func (r *PublicRouter) ProvideHandlers() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(actions.NewRegisterAction),
		flow.NewProvider(actions.NewLoginAction),
	}
}

func (r *PublicRouter) RegisterSubRouters() bool {
	return false
}
