package dialect

import (
	"fmt"

	"api/providers/db"
)

type mssql struct {
	db db.Store
}

func init() {
	RegisterDialect("mssql", &mssql{})
}

func (mssql) Name() string {
	return "mssql"
}

func (c *mssql) SetDB(db db.Store) {
	c.db = db
}

func (c *mssql) DB() db.Store {
	return c.db
}

func (c *mssql) CurrentDatabase() (name string) {
	c.db.QueryRow("SELECT DB_NAME() AS [Current Database]").Scan(&name)
	return
}

func (c *mssql) HasTable(tableName string) bool {
	var count int
	currentDatabase := c.CurrentDatabase()
	c.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = ? AND table_catalog = ?", tableName, currentDatabase).Scan(&count)
	return count > 0
}

func (c *mssql) MigrationExists(version string, tableName string) (bool, error) {
	return false, fmt.Errorf("Function `MigrationExists` not implemented for dialect `%s`", c.Name())
}

func (c *mssql) CreateMigrationTable(tableName string) error {
	return fmt.Errorf("Function `CreateMigrationTable` not implemented for dialect `%s`", c.Name())
}

func (c *mssql) CountRecords(tableName string) (int, error) {
	return 0, fmt.Errorf("Function `CountRecords` not implemented for dialect `%s`", c.Name())
}

func (c *mssql) SaveMigration(tableName string, version string, name string) error {
	return fmt.Errorf("Function `SaveMigration` not implemented for dialect `%s`", c.Name())
}

func (c *mssql) RemoveMigration(tableName string, version string) error {
	return fmt.Errorf("Function `RemoveMigration` not implemented for dialect `%s`", c.Name())
}
