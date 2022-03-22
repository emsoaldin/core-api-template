package routers

import (
	"api/providers/jwt"

	"github.com/go-flow/flow/v2"
)

// Router
type Router struct {
	auth jwt.TokenAuth
}

func NewRouter(auth jwt.TokenAuth) *Router {
	return &Router{
		auth: auth,
	}
}

// Path defined http path for router
func (r *Router) Path() string {
	return "/users"
}

// Middlewares provides list of middlewares used by the router
func (r *Router) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{
		r.auth.AuthorizeRequest("Authorized"),
		r.auth.AuthorizeRequest("Admin"),
	}
}

func (r *Router) ProvideHandlers() []flow.Provider {
	return []flow.Provider{}
}

func (r *Router) RegisterSubRouters() bool {
	return false
}
