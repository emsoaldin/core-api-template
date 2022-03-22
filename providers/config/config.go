package config

import (
	"log"
	"os"
)

//AppConfig holds all application configuration
type AppConfig interface {
	// Env returns execution environment configuration
	Env() string

	// LogLevel returns log level configuration
	LogLevel() string

	// Addr returns http serving listen address
	Addr() string

	// RSAPrivateKey returns string content for RSA key
	RSAPrivateKey() string

	// RSAPublicKey returns string content for RSA key
	RSAPublicKey() string

	// RSAKeyPassword returns password for RSA key
	RSAKeyPassword() string
}

// New creates new Configuration object
func New() AppConfig {
	privateKeyPath := mustGetEnv("RSA_PRIVATE_KEY")
	publicKeyPath := mustGetEnv("RSA_PUBLIC_KEY")
	privateKeyPwd := getEnv("RSA_PRIVATE_KEY_PASSWORD", "")

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	return &config{
		env:            getEnv("ENV", "development"),
		logLevel:       getEnv("LOG_LEVEL", "error"),
		addr:           getEnv("ADDR", ""),
		rsaPrivateKey:  string(privateKey),
		rsaPublicKey:   string(publicKey),
		rsaKeyPassword: privateKeyPwd,
	}
}

type config struct {
	env            string
	logLevel       string
	addr           string
	rsaPrivateKey  string
	rsaPublicKey   string
	rsaKeyPassword string
}

// Env returns execution environment configuration
func (c *config) Env() string {
	return c.env
}

// LogLevel returns log level configuration
func (c *config) LogLevel() string {
	return c.logLevel
}

// Addr returns http serving listen address
func (c *config) Addr() string {
	return c.addr
}

// RSAPrivateKey returns string content for RSA key
func (c *config) RSAPrivateKey() string {
	return c.rsaPrivateKey
}

// RSAPublicKey returns string content for RSA key
func (c *config) RSAPublicKey() string {
	return c.rsaPublicKey
}

// RSAKeyPassword returns password for RSA key
func (c *config) RSAKeyPassword() string {
	return c.rsaKeyPassword
}

// getEnv returns value for given key from environment
// if key is not present in environment it returns defaultValue
func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if len(v) > 0 {
		return v
	}
	return defaultValue
}

// mustGetEnv returns value for given key from environment
// if key is not present in environment function will panic
func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		log.Fatalf(" variable `%s` is not present in ENVIRONMENT", key)
	}
	return v
}
