package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"api/providers/config"
	"api/providers/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-flow/flow/v2"
)

// ErrForbiddden is returned when user does not have required claims
var ErrForbiddden = errors.New("forbidden")

var ErrUnathorized = errors.New("unauthorized")

// TokenAuth interface
type TokenAuth interface {
	// TokenAuth returns service implementation signature
	TokenAuth() string

	JWKS() *JWKS

	GenerateAccessToken(userID uint64, scope ...string) (string, error)

	VerifyAccessToken(token string, scope ...string) (uint64, string, error)

	GenerateRefreshToken(token string) (string, error)

	VerifyRefreshToken(tokenString string) (string, error)

	GenerateAPIKey(userID uint64) (string, error)

	VerifyAPIKey(tokenString string) (uint64, error)

	AuthorizeRequest(roles ...string) flow.MiddlewareHandlerFunc

	AuthorizeAPIKey() flow.MiddlewareHandlerFunc

	// RequestUserID returns userId from request context
	// if user isid is not found then unathorized error is returned
	RequestUserID(r *http.Request) (uint64, error)

	// RequestUserClaims returns authorization claims from request context
	// if claims are not found then unathorized error is returned
	RequestUserClaims(r *http.Request) ([]string, error)
}

func NewAuth(cfg config.AppConfig, logger log.Logger) TokenAuth {
	var privateKey *rsa.PrivateKey
	var err error

	if cfg.RSAKeyPassword() != "" {
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEMWithPassword([]byte(cfg.RSAPrivateKey()), cfg.RSAKeyPassword())
	} else {
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.RSAPrivateKey()))
	}

	if err != nil {
		logger.Fatal(err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.RSAPublicKey()))

	if err != nil {
		logger.Fatal(err)
	}

	privateKey.PublicKey = *publicKey

	modulus := publicKey.N
	exponent := big.NewInt(int64(publicKey.E))

	jwk := new(JWK)
	jwk.Alg = "RS256"
	jwk.Use = "sig"
	jwk.Kid = "base"
	jwk.Kty = "RSA"
	jwk.N = base64.RawURLEncoding.EncodeToString(modulus.Bytes())
	jwk.E = base64.RawURLEncoding.EncodeToString(exponent.Bytes())

	block, _ := pem.Decode([]byte(cfg.RSAPublicKey()))
	jwk.X5C = []string{
		base64.StdEncoding.EncodeToString(block.Bytes),
	}

	jwks := new(JWKS)
	jwks.Keys = append(jwks.Keys, jwk)

	return &jwtTokenAuth{
		privateKey: privateKey,
		jwks:       jwks,
	}
}

type jwtTokenAuth struct {
	privateKey *rsa.PrivateKey
	jwks       *JWKS
}

func (jwtTokenAuth) TokenAuth() string {
	return "jwtTokenAuth"
}

func (svc *jwtTokenAuth) JWKS() *JWKS {
	return svc.jwks
}

func (svc *jwtTokenAuth) GenerateAccessToken(userID uint64, scope ...string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := token.Claims.(jwt.MapClaims)

	token.Header["kid"] = "base"
	token.Header["typ"] = "JWT"

	now := time.Now().UTC().Unix()

	claims["uid"] = userID
	claims["scope"] = strings.Join(scope, " ")
	expIn := time.Minute * time.Duration(16)
	claims["exp"] = now + int64(expIn.Seconds())
	claims["iat"] = now

	return token.SignedString(svc.privateKey)
}

func (svc *jwtTokenAuth) VerifyAccessToken(accessToken string, claims ...string) (uint64, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != token.Method {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return &svc.privateKey.PublicKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	if token.Header["typ"] != "JWT" {
		return 0, "", fmt.Errorf("not an access token")
	}

	if c, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint64(c["uid"].(float64))
		scope := c["scope"].(string)
		err = svc.verifyClaims(scope, claims...)
		return userID, scope, err
	}

	return 0, "", fmt.Errorf("invalid token")
}

func (svc *jwtTokenAuth) GenerateRefreshToken(token string) (string, error) {
	refreshToken := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := refreshToken.Claims.(jwt.MapClaims)

	refreshToken.Header["kid"] = "base"
	refreshToken.Header["typ"] = "RFRSH"

	now := time.Now().UTC().Unix()
	claims["token"] = token
	expIn := time.Hour * time.Duration(24) * 30 // 30 Days
	claims["exp"] = now + int64(expIn.Seconds())
	claims["iat"] = now

	return refreshToken.SignedString(svc.privateKey)
}

func (svc *jwtTokenAuth) VerifyRefreshToken(tokenString string) (token string, err error) {
	refreshToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != token.Method {
			return nil, fmt.Errorf("invalid signing algorithm")
		}

		return &svc.privateKey.PublicKey, nil
	})

	if refreshToken == nil {
		err = errors.New("invalid token")
		return
	}

	if refreshToken.Header["typ"] != "RFRSH" {
		err = fmt.Errorf("not a refresh token")
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		token = claims["token"].(string)
	}

	return
}

func (svc *jwtTokenAuth) GenerateAPIKey(userID uint64) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := token.Claims.(jwt.MapClaims)

	token.Header["kid"] = "base"
	token.Header["typ"] = "APIKEY"

	now := time.Now().UTC().Unix()

	claims["uid"] = userID
	// Yes, api key expires in 100 years
	expIn := time.Hour * time.Duration(876000)
	claims["exp"] = now + int64(expIn.Seconds())
	claims["iat"] = now

	return token.SignedString(svc.privateKey)
}

func (svc *jwtTokenAuth) VerifyAPIKey(apiKey string) (uint64, error) {
	token, err := jwt.Parse(apiKey, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != token.Method {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return &svc.privateKey.PublicKey, nil
	})

	if err != nil {
		return 0, err
	}

	if token.Header["typ"] != "APIKEY" {
		return 0, fmt.Errorf("not an api key")
	}

	if c, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint64(c["uid"].(float64))
		return userID, err
	}

	return 0, fmt.Errorf("invalid api key")
}

func (svc *jwtTokenAuth) AuthorizeRequest(claims ...string) flow.MiddlewareHandlerFunc {
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {
			token := svc.findAuthorizationToken(r)

			id, scope, err := svc.VerifyAccessToken(token, claims...)

			if err != nil {
				if errors.Is(err, ErrForbiddden) {
					return flow.ResponseError(http.StatusForbidden, err)
				}
				return flow.ResponseError(http.StatusUnauthorized, err)
			}

			// add id and scope claims to request
			ctx := r.Context()
			ctx = NewIDClaimContext(ctx, id)
			ctx = NewScopeClaimContext(ctx, scope)
			r = r.WithContext(ctx)

			return next(w, r)
		}
	}
}

func (svc *jwtTokenAuth) AuthorizeAPIKey() flow.MiddlewareHandlerFunc {
	return func(next flow.MiddlewareFunc) flow.MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) flow.Response {
			token := svc.findApiKey(r)

			id, err := svc.VerifyAPIKey(token)

			if err != nil {
				if errors.Is(err, ErrForbiddden) {
					return flow.ResponseError(http.StatusForbidden, err)
				}
				return flow.ResponseError(http.StatusUnauthorized, err)
			}

			// add id and scope claims to request
			ctx := r.Context()
			ctx = NewIDClaimContext(ctx, id)
			r = r.WithContext(ctx)

			return next(w, r)
		}
	}
}

// RequestUserID returns userId from request context
// if user isid is not found then unathorized error is returned
func (svc *jwtTokenAuth) RequestUserID(r *http.Request) (uint64, error) {
	userId, ok := IDClaimFromContext(r.Context())
	if !ok {
		return 0, ErrUnathorized
	}
	return userId, nil
}

// RequestUserClaims returns authorization claims from request context
// if claims are not found then unauthorized error is returned
func (svc *jwtTokenAuth) RequestUserClaims(r *http.Request) ([]string, error) {
	claims, ok := ScopeClaimFromContext(r.Context())
	if !ok {
		return []string{}, ErrUnathorized
	}
	return strings.Split(claims, ","), nil
}

func (svc *jwtTokenAuth) findAuthorizationToken(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func (svc *jwtTokenAuth) findApiKey(r *http.Request) string {
	// Get api key from X-API-KEY header.
	apiKey := r.Header.Get("X-API-KEY")
	return apiKey
}

func (svc *jwtTokenAuth) verifyClaims(scope interface{}, claims ...string) error {
	if len(claims) > 0 {
		// cast scope to string
		scopeStr, ok := scope.(string)
		if !ok {
			return ErrForbiddden
		}

		match := false
		for _, claim := range claims {
			if strings.Contains(scopeStr, claim) {
				match = true
				break
			}
		}
		if !match {
			return ErrForbiddden
		}
	}

	return nil
}
