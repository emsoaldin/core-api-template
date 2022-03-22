package dialect

import (
	"fmt"
	"reflect"

	"api/providers/db"
)

// Dialect defines set of methods needed
// by migrator to interact with different database engines
//
// It contains behaviors that differ across SQL database
type Dialect interface {
	// Name get dialect's name
	Name() string

	// SetDB set db for dialect
	SetDB(db db.Store)

	// DB gets dialect db
	DB() db.Store

	// HasTable check has table or not
	HasTable(tableName string) bool

	// CurrentDatabase return current database name
	CurrentDatabase() string

	// MigrationExists  checks if migration version exists in database table
	MigrationExists(version string, tableName string) (bool, error)

	// CreateMigrationTable creates migrations table in database
	CreateMigrationTable(tableName string) error

	// CountRecords retunrs number of rows in provided table
	CountRecords(tableName string) (int, error)

	// SaveMigration stores migration version and name in database
	SaveMigration(tableName string, version string, name string) error

	// RemoveMigration deletes migration version from database table
	RemoveMigration(tableName string, version string) error
}

var dialectsMap = map[string]Dialect{}

// New creates dialect instance for given dialect name and db connection
func New(name string, db db.Store) Dialect {
	if value, ok := dialectsMap[name]; ok {
		dialect := reflect.New(reflect.TypeOf(value).Elem()).Interface().(Dialect)
		dialect.SetDB(db)
		return dialect
	}

	panic(fmt.Sprintf("dialect `%s` is not supported", name))
}

// RegisterDialect register new dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}
