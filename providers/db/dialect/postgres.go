package dialect

import (
	"fmt"

	"api/providers/db"
)

type postgres struct {
	db db.Store
}

func init() {
	RegisterDialect("postgres", &postgres{})
	RegisterDialect("cloudsqlpostgres", &postgres{})
}

func (postgres) Name() string {
	return "postgres"
}

func (c *postgres) SetDB(db db.Store) {
	c.db = db
}

func (c *postgres) DB() db.Store {
	return c.db
}

func (c *postgres) CurrentDatabase() (name string) {
	c.db.QueryRow("SELECT CURRENT_DATABASE()").Scan(&name)
	return
}

func (c *postgres) HasTable(tableName string) bool {
	var count int
	c.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = $1 AND table_type = 'BASE TABLE' AND table_schema = CURRENT_SCHEMA()", tableName).Scan(&count)
	return count > 0
}

func (c *postgres) MigrationExists(version string, tableName string) (bool, error) {
	return false, fmt.Errorf("Function `MigrationExists` not implemented for dialect `%s`", c.Name())
}

func (c *postgres) CreateMigrationTable(tableName string) error {
	return fmt.Errorf("Function `CreateMigrationTable` not implemented for dialect `%s`", c.Name())
}

func (c *postgres) CountRecords(tableName string) (int, error) {
	return 0, fmt.Errorf("Function `CountRecords` not implemented for dialect `%s`", c.Name())
}

func (c *postgres) SaveMigration(tableName string, version string, name string) error {
	return fmt.Errorf("Function `SaveMigration` not implemented for dialect `%s`", c.Name())
}

func (c *postgres) RemoveMigration(tableName string, version string) error {
	return fmt.Errorf("Function `RemoveMigration` not implemented for dialect `%s`", c.Name())
}
