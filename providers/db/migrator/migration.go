package migrator

import (
	"fmt"

	"api/providers/db"
)

// Migration handles data for a given database migration
type Migration struct {

	//Version of the migration
	//
	// for example 2`018051570142`
	Version string

	// Name of the migration
	//
	// for exmample `initial`
	Name string

	// Direction of the migration
	//
	// for example `up`
	Direction string

	// Content of migration
	//
	// Content holds migration SQL query
	Content string
}

// Run executes the migration.
//
// Returns error if Content is not defined,
// and returns result from SQL execution (error)
func (m Migration) Run(conn db.Store) error {
	if m.Content == "" {
		return fmt.Errorf("migration runner not defined for %s.%s.%s", m.Version, m.Name, m.Direction)
	}

	// get DB transaction
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	//execute transaction
	_, err = tx.Exec(m.Content)
	if err != nil {
		// rollback
		err1 := tx.Rollback()
		// return rollback error
		if err1 != nil {
			return err1
		}
		// return tx execution error
		return err
	}
	// return tx commit result - success will return nil
	return tx.Commit()
}

// Migrations collection
type Migrations []Migration

// Len returns number of migrations
func (m Migrations) Len() int {
	return len(m)
}

// Less checks if migration version at index i is lesss than migration at index j
func (m Migrations) Less(i, j int) bool {
	return m[i].Version < m[j].Version
}

// Swap migrations at index i and j
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
