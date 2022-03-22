package middlewares

import (
	"net/http"

	"github.com/go-flow/flow/v2"
)

// Sample handles sample logic
func Sample() flow.MiddlewareHandlerFunc {
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {
			//code that will be executed before next middleware in pipeline

			// invoke next middleware
			resp := next(w, r)

			// code that will be executed after next middleware

			// return retult error from previous middleware
			return resp
		}
	}
}
