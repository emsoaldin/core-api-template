package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
)

func New() Store {
	//root:root@(db:3306)/core_api?multiStatements=true&readTimeout=1800s&charset=utf8mb4&parseTime=True&loc=Local
	var (
		dbDialect = getEnv("DB_DIALECT", "mysql")
		dbUser    = mustGetEnv("DB_USER")     // e.g. 'my-db-user'
		dbPwd     = mustGetEnv("DB_PASS")     // e.g. 'my-db-password'
		dbHost    = mustGetEnv("DB_HOST")     // e.g. '/cloudsql/project:region:instance'
		dbPort    = getEnv("DB_PORT", "3306") // e.g. '3306'
		dbName    = mustGetEnv("DB_NAME")     // e.g. 'my-database'
		dbParams  = getEnv("DB_PARAMS", "")   // e.g. 'parseTime=true'

		maxIdleConns    = getEnvInt("DB_MAX_IDLE_CONNS", 10)
		maxOpenConns    = getEnvInt("DB_MAX_OPEN_CONNS", 100)
		connMaxLifetime = getEnvInt("DB_MAX_LIFETIME", 30)
	)

	dbURI := createUri(dbDialect, dbUser, dbPwd, dbHost, dbPort, dbName, dbParams)

	if dbDialect == "cloudsql" {
		dbDialect = "mysql"
	}

	dbPool, err := sql.Open(dbDialect, dbURI)
	if err != nil {
		panic(fmt.Errorf("unable to open DB connection: %w", err))
	}

	// SetMaxIdleConns sets maximum number of connections in the idle connection pool
	dbPool.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	dbPool.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	duration := time.Minute * time.Duration(connMaxLifetime)
	dbPool.SetConnMaxLifetime(duration)

	// ping DB
	if err := dbPool.Ping(); err != nil {
		panic(fmt.Errorf("unable to ping Database: %w", err))
	}

	return dbPool

}

func createUri(dbDialect, dbUser, dbPwd, dbHost, dbPort, dbName, dbParams string) string {
	if dbDialect == "cloudsql" {
		return fmt.Sprintf("%s:%s@unix(%s)/%s?%s", dbUser, dbPwd, dbHost, dbName, dbParams)

	}

	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", dbUser, dbPwd, dbHost, dbPort, dbName, dbParams)
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
		panic(fmt.Errorf("variable `%s` is not present in ENVIRONMENT", key))
	}
	return v
}

// getEnvInt returns integer value for given key from environment
// if key is not present in environment it returns defaultValue
// if key cannot be parsed to integer function will panic
func getEnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defaultValue
	}

	valInteger, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Errorf(" variable `%s` cannot be parsed to INTEGER", key))

	}

	return valInteger
}
