package dialect

import (
	"fmt"

	"api/providers/db"
)

type sqlite3 struct {
	db db.Store
}

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

func (sqlite3) Name() string {
	return "sqlite3"
}

func (c *sqlite3) SetDB(db db.Store) {
	c.db = db
}

func (c *sqlite3) DB() db.Store {
	return c.db
}

func (c *sqlite3) CurrentDatabase() (name string) {
	var (
		ifaces   = make([]interface{}, 3)
		pointers = make([]*string, 3)
		i        int
	)
	for i = 0; i < 3; i++ {
		ifaces[i] = &pointers[i]
	}
	if err := c.db.QueryRow("PRAGMA database_list").Scan(ifaces...); err != nil {
		return
	}
	if pointers[1] != nil {
		name = *pointers[1]
	}
	return
}

func (c *sqlite3) HasTable(tableName string) bool {
	var count int
	c.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
	return count > 0
}

func (c *sqlite3) MigrationExists(version string, tableName string) (bool, error) {
	return false, fmt.Errorf("Function `MigrationExists` not implemented for dialect `%s`", c.Name())
}

func (c *sqlite3) CreateMigrationTable(tableName string) error {
	return fmt.Errorf("Function `CreateMigrationTable` not implemented for dialect `%s`", c.Name())
}

func (c *sqlite3) CountRecords(tableName string) (int, error) {
	return 0, fmt.Errorf("Function `CountRecords` not implemented for dialect `%s`", c.Name())
}

func (c *sqlite3) SaveMigration(tableName string, version string, name string) error {
	return fmt.Errorf("Function `SaveMigration` not implemented for dialect `%s`", c.Name())
}

func (c *sqlite3) RemoveMigration(tableName string, version string) error {
	return fmt.Errorf("Function `RemoveMigration` not implemented for dialect `%s`", c.Name())
}
