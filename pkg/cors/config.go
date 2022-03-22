package cors

import (
	"errors"
	"net/http"
	"strings"
)

type cors struct {
	allowAllOrigins  bool
	allowCredentials bool
	allowOriginFunc  func(string) bool
	allowOrigins     []string
	normalHeaders    http.Header
	preflightHeaders http.Header
	wildcardOrigins  [][]string
}

var (
	// DefaultSchemas holds default allowed schemas
	DefaultSchemas = []string{
		"http://",
		"https://",
	}

	// ExtensionSchemas holds default browser extension schemas
	ExtensionSchemas = []string{
		"chrome-extension://",
		"safari-extension://",
		"moz-extension://",
		"ms-browser-extension://",
	}

	// FileSchemas hodls default file schemas
	FileSchemas = []string{
		"file://",
	}
	// WebSocketSchemas holds default WebSocket Schemas
	WebSocketSchemas = []string{
		"ws://",
		"wss://",
	}

	// ErrForbidden is returned on invalid origins
	ErrForbidden = errors.New("Forbidden")
)

func newCors(config Config) *cors {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}

	return &cors{
		allowOriginFunc:  config.AllowOriginFunc,
		allowAllOrigins:  config.AllowAllOrigins,
		allowCredentials: config.AllowCredentials,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalHeaders(config),
		preflightHeaders: generatePreflightHeaders(config),
		wildcardOrigins:  config.parseWildcardRules(),
	}
}

func (cors *cors) applyCors(w http.ResponseWriter, r *http.Request) error {
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return nil
	}
	host := r.Header.Get("Host")
	if origin == "http://"+host || origin == "https://"+host {
		// request is not a CORS request but have origin header.
		// for example, use fetch api
		return nil
	}

	if !cors.validateOrigin(origin) {
		return ErrForbidden
	}

	if r.Method == "OPTIONS" {
		cors.handlePreflight(w)
	} else {
		cors.handleNormal(w)
	}

	if !cors.allowAllOrigins {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	return nil
}

func (cors *cors) validateWildcardOrigin(origin string) bool {
	for _, w := range cors.wildcardOrigins {
		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
			return true
		}
		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
			return true
		}
		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
			return true
		}
	}

	return false
}

func (cors *cors) validateOrigin(origin string) bool {
	if cors.allowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if len(cors.wildcardOrigins) > 0 && cors.validateWildcardOrigin(origin) {
		return true
	}
	if cors.allowOriginFunc != nil {
		return cors.allowOriginFunc(origin)
	}
	return false
}

func (cors *cors) handlePreflight(w http.ResponseWriter) {
	header := w.Header()
	for key, value := range cors.preflightHeaders {
		header[key] = value
	}
}

func (cors *cors) handleNormal(w http.ResponseWriter) {
	header := w.Header()
	for key, value := range cors.normalHeaders {
		header[key] = value
	}
}
