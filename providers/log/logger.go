package log

import "os"

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger interface
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// New creates new Logger instance
func New() Logger {
	return newZapLogger(Configuration{
		JSONFormat: true,
		Level:      getEnv("LOG_LEVEL", "debug"),
		Output:     os.Stdout,
	})
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
