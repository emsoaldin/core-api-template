package cors

import (
	"errors"
	"net/http"

	"github.com/go-flow/flow/v2"
)

// Middleware returns the cors middleware with default configuration.
func Middleware() flow.MiddlewareHandlerFunc {
	config := DefaultConfig()
	config.AllowAllOrigins = true
	return MiddlewareWithConfig(config)
}

// MiddlewareWithConfig returns the cors middleware with user-defined custom configuration.
func MiddlewareWithConfig(config Config) flow.MiddlewareHandlerFunc {
	cors := newCors(config)
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {
			if err := cors.applyCors(w, r); err != nil {
				if errors.Is(ErrForbidden, err) {
					return flow.ResponseError(http.StatusForbidden, err)
				}
				return flow.ResponseError(http.StatusInternalServerError, err)
			}

			return next(w, r)
		}
	}
}
