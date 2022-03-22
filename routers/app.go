package routers

import (
	"api/actions"
	"api/pkg/cors"
	"api/providers/config"
	"api/providers/db"
	"api/providers/log"
	"api/version"
	"fmt"
	"net/http"

	"github.com/go-flow/flow/v2"
)

// Router
type Router struct {
	config config.AppConfig
	logger log.Logger
	store  db.Store
}

func NewRouter(config config.AppConfig, logger log.Logger, store db.Store) *Router {
	return &Router{
		config: config,
		logger: logger,
		store:  store,
	}
}

// Path defined http path for router
func (r *Router) Path() string {
	return "/"
}

func (r *Router) PanicRecover(next flow.MiddlewareFunc) flow.MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) flow.Response {
		defer func() {
			if err := recover(); err != nil {
				resp := flow.ResponseError(http.StatusInternalServerError, fmt.Errorf("panic: %v", err))
				if e := resp.Handle(w, r); e != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprintf("Response error: %v", e)))
				}
			}
		}()

		return next(w, r)
	}
}

// Middlewares provides list of middlewares used by the router
func (r *Router) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{
		r.PanicRecover,
		log.MiddlewareWithFields(r.logger, log.Fields{"version": version.Build}),
		cors.Middleware(),
		db.RequestTxMiddleware(r.store),
	}
}

func (r *Router) ProvideHandlers() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(actions.NewIndexAction),
		flow.NewProvider(actions.NewHealthAction),
		flow.NewProvider(actions.NewSwaggerAction),
	}
}

func (r *Router) RegisterSubRouters() bool {
	return true
}
