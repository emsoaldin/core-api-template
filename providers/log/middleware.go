package log

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-flow/flow/v2"
	"github.com/rs/xid"
)

// Middleware creates request logger for given logger instance
func Middleware(logger Logger, ignore ...string) flow.MiddlewareHandlerFunc {
	return MiddlewareWithFields(logger, Fields{}, ignore...)
}

// MiddlewareWithFields creates request logger for given logger instance and custom fileds
func MiddlewareWithFields(logger Logger, f Fields, ignore ...string) flow.MiddlewareHandlerFunc {
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {
			start := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = r.Header.Get("X-Amzn-Trace-Id")
			}

			if requestID == "" {
				// generate new RequestID
				guid := xid.New()
				requestID = guid.String()
				// add requestID to header
				r.Header.Add("X-Request-ID", requestID)
			}

			fields := make(map[string]interface{})

			for k, v := range f {
				fields[k] = v
			}

			// add requestID to response header too
			w.Header().Add("X-Request-ID", requestID)

			fields["request_id"] = requestID

			// create request logger instance
			rl := logger.WithFields(fields)

			// put request logger to request context for later use
			ctx := r.Context()
			ctx = NewContext(ctx, rl)

			r = r.WithContext(ctx)

			// invoke next middleware
			resp := next(w, r)

			// check if current request logging should be skipped
			ignoreList := strings.Join(ignore, ",")
			if strings.Contains(ignoreList, r.URL.Path) {
				return resp
			}

			rl.WithFields(Fields{
				"status":   resp.Status(),
				"method":   r.Method,
				"path":     r.URL.String(),
				"duration": time.Since(start).String(),
			}).Info("request-log")

			// return response from previous middleware
			return resp
		}
	}
}
