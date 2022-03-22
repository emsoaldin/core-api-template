package jwt

//JWK model
type JWK struct {
	Alg string   `json:"alg,omitempty"`
	E   string   `json:"e,omitempty"`
	Kid string   `json:"kid,omitempty"`
	Kty string   `json:"kty,omitempty"`
	N   string   `json:"n,omitempty"`
	Use string   `json:"use,omitempty"`
	X5C []string `json:"x5c,omitempty"`
}

// JWKS model
type JWKS struct {
	Keys []*JWK `json:"keys,omitempty"`
}
