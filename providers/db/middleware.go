package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-flow/flow/v2"
)

// RequestTxMiddleware creates request scoped transaction
func RequestTxMiddleware(db Store) flow.MiddlewareHandlerFunc {
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {

			// create database context aware transaction
			tx, err := db.BeginTx(r.Context(), nil)

			if err != nil {
				return flow.ResponseError(http.StatusInternalServerError, fmt.Errorf("unable to create request scoped db.Tx; %w", err))
			}

			// store tx object in context
			ctx := r.Context()
			ctx = NewTxContext(ctx, tx)
			r = r.WithContext(ctx)

			// invoke next middleware
			res := next(w, r)

			if res.Status() < 200 || res.Status() >= 400 {
				// response status is not in allowed range
				if err := tx.Rollback(); err != nil {
					return flow.ResponseError(http.StatusInternalServerError, fmt.Errorf("unable to Rollback request scoped db.Tx; %w", err))
				}

			}

			if err = tx.Commit(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				return flow.ResponseError(http.StatusInternalServerError, fmt.Errorf("unable to Commit request scoped db.Tx; %w", err))
			}

			return res
		}
	}
}
