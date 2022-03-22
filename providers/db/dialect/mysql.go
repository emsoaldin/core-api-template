package dialect

import (
	"fmt"

	"api/providers/db"
)

type mysql struct {
	db db.Store
}

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (mysql) Name() string {
	return "mysql"
}

func (c *mysql) SetDB(db db.Store) {
	c.db = db
}

func (c *mysql) DB() db.Store {
	return c.db
}

func (c *mysql) CurrentDatabase() (name string) {
	c.db.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

func (c *mysql) HasTable(tableName string) bool {
	var count int
	currentDatabase := c.CurrentDatabase()
	c.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, tableName).Scan(&count)
	return count > 0
}

func (c *mysql) MigrationExists(version string, tableName string) (bool, error) {
	var count int
	qStr := fmt.Sprintf("SELECT COUNT(*) FROM %s where version = ?", tableName)
	err := c.db.QueryRow(qStr, version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *mysql) CreateMigrationTable(tableName string) error {
	q := fmt.Sprintf(`CREATE TABLE %s ( 
		version NVARCHAR(14) NOT NULL, 
		name NVARCHAR(255) NULL, 
		UNIQUE INDEX  schema_version_idx (version ASC));`, tableName)
	_, err := c.db.Exec(q)
	return err
}

func (c *mysql) CountRecords(tableName string) (int, error) {
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int
	err := c.db.QueryRow(q).Scan(&count)
	return count, err
}

func (c *mysql) SaveMigration(tableName string, version string, name string) error {
	q := fmt.Sprintf("INSERT INTO %s (version, name) VALUES (?,?)", tableName)
	_, err := c.db.Exec(q, version, name)
	return err
}

func (c *mysql) RemoveMigration(tableName string, version string) error {
	q := fmt.Sprintf("delete from %s where version = ?", tableName)
	_, err := c.db.Exec(q, version)
	return err
}
