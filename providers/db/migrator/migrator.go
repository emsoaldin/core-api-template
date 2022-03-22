package migrator

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"api/providers/db"
	"api/providers/db/dialect"
)

var (
	// migrationRegEx defines regular expression for migration files
	//
	// example file `2018051570142_initial.up.sql`
	migrationRegEx = regexp.MustCompile(`(\d+)_([^\.]+)(\.[a-z]+)?\.(up|down)\.(sql)`)

	// MigrationsTableName holds default table name for migrations
	//
	// Default value is "schema_migration"
	MigrationsTableName = "schema_migration"
)

// newMigrator returns a new "blank" migrator.
//
// a blank Migrator should be only used as
// basis for a new type of migration system
func newMigrator(dialectName string, conn db.Store) Migrator {
	m := Migrator{
		dialect: dialect.New(dialectName, conn),
		Migrations: map[string]Migrations{
			"up":   {},
			"down": {},
		},
	}
	// create migration schema
	err := m.createMigrationSchema()
	if err != nil {
		panic(err)
	}

	return m
}

// Migrator forms the basis of all migration systems.Migrator
// it does the actual heavy lifting of running migrations
type Migrator struct {
	dialect      dialect.Dialect
	Migrations   map[string]Migrations
	lastExecTime string
}

// Up runs pending `up` migrations and applies them to the database
func (m Migrator) Up() error {
	return m.exec(func() error {

		migrations := m.Migrations["up"]
		sort.Sort(migrations)
		for _, migration := range migrations {
			exists, err := m.migrationExists(migration)
			if err != nil {
				return err
			}

			if exists {
				continue // migration is already
			}

			err = migration.Run(m.dialect.DB())
			if err != nil {
				return err
			}

			err = m.dialect.SaveMigration(MigrationsTableName, migration.Version, migration.Name)
			if err != nil {
				return err
			}

			fmt.Printf("> %s\n", migration.Name)
		}
		return nil
	})
}

// Down runs pending `down` migrations and
//rolls back the database by the specified number of steps
func (m Migrator) Down(step int) error {
	return m.exec(func() error {
		count, err := m.executedMigrationsCount()
		if err != nil {
			return err
		}

		migrations := m.Migrations["down"]
		sort.Sort(sort.Reverse(migrations))
		//skip all executed migrations
		if len(migrations) > count {
			migrations = migrations[len(migrations)-count:]
		}

		// run only required steps
		if step > 0 && len(migrations) >= step {
			migrations = migrations[:step]
		}

		for _, migration := range migrations {
			exists, err := m.migrationExists(migration)
			if err != nil {
				return nil
			}

			if !exists {
				return fmt.Errorf("problem checking for migration version `%s`", migration.Version)
			}

			err = migration.Run(m.dialect.DB())
			if err != nil {
				return err
			}

			err = m.dialect.RemoveMigration(MigrationsTableName, migration.Version)
			if err != nil {
				return err
			}

			fmt.Printf("< %s\n", migration.Name)
		}

		return nil
	})
}

// Status prints out the status of applied/Pending migrations
func (m Migrator) Status() error {
	return m.exec(func() error {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "Version\t\tName\t\tStatus")
		for _, migration := range m.Migrations["up"] {
			exists, err := m.migrationExists(migration)
			if err != nil {
				return err
			}
			status := "Pending"
			if exists {
				status = "Applied"
			}

			fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\t\n", migration.Version, migration.Name, status)
		}
		return w.Flush()
	})
}

// Reset executes all `down` migrations followed by the `up` migrations
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return err
	}
	return m.Up()
}

// LastActionExecutionTime returns execution time of last migrator action execution
func (m Migrator) LastActionExecutionTime() string {
	return m.lastExecTime
}

func (m Migrator) hasMigrationSchema() bool {
	return m.dialect.HasTable(MigrationsTableName)
}

func (m Migrator) createMigrationSchema() error {
	if m.hasMigrationSchema() {
		return nil // migration exists
	}

	return m.dialect.CreateMigrationTable(MigrationsTableName)
}

func (m Migrator) executedMigrationsCount() (int, error) {
	return m.dialect.CountRecords(MigrationsTableName)
}

// Exists checks if migration exists in DB
func (m Migrator) migrationExists(migration Migration) (bool, error) {
	return m.dialect.MigrationExists(migration.Version, MigrationsTableName)
}

// exec internal helper execution function which prints
// execusion time of passed function
func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer func() {
		diff := time.Now().Sub(now).Seconds()
		if diff > 60 {
			m.lastExecTime = fmt.Sprintf("\n%.4f minutes \n", diff/60)
		} else {
			m.lastExecTime = fmt.Sprintf("\n%.4f seconds \n", diff)
		}

	}()

	return fn()
}
