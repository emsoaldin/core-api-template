package jwt

import (
	"context"
)

type IDClaimKey struct{}

type ScopeClaimKey struct{}

// IDClaimFromContext returns ID Claim value within given context
func IDClaimFromContext(ctx context.Context) (uint64, bool) {
	tx, ok := ctx.Value(IDClaimKey{}).(uint64)
	return tx, ok
}

// NewIDClaimContext creates context with IDClaim value
func NewIDClaimContext(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, IDClaimKey{}, id)
}

// ScopeClaimFromContext returns Scope claims value within given context
func ScopeClaimFromContext(ctx context.Context) (string, bool) {
	tx, ok := ctx.Value(ScopeClaimKey{}).(string)
	return tx, ok
}

// NewScopeClaimContext creates context with Scope Claims value
func NewScopeClaimContext(ctx context.Context, scope string) context.Context {
	return context.WithValue(ctx, ScopeClaimKey{}, scope)
}
