package migrator

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"api/providers/db"
)

// FSMigrator is Migrator implementation for SQL
// files stored on file system using embed.FS
type FSMigrator struct {
	Migrator
	FS embed.FS
}

// NewFSMigrator - creates Migrations for embed.FS
func NewFSMigrator(fs embed.FS, dialect string, store db.Store) FSMigrator {
	fm := FSMigrator{
		Migrator: newMigrator(dialect, store),
		FS:       fs,
	}

	err := fm.loadMigrations()
	if err != nil {
		panic(err)
	}
	return fm
}

func (fsm *FSMigrator) loadMigrations() error {

	files, err := fsm.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("unable to read embed.FS Directory. Error: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()
		matches := migrationRegEx.FindAllStringSubmatch(fileName, -1)
		if matches == nil || len(matches) == 0 {
			return fmt.Errorf("file %s does not match migration file pattern", fileName)
		}
		raw, err := fsm.FS.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("unable to read %s File. Error: %w", fileName, err)
		}
		content := string(raw)
		temp := template.Must(template.New("sql").Parse(content))
		var buff bytes.Buffer
		err = temp.Execute(&buff, nil)
		if err != nil {
			return fmt.Errorf("unable to parse %s file. Error: %w", fileName, err)
		}

		content = buff.String()

		match := matches[0]
		dir := match[4]

		migration := Migration{
			Version:   match[1],
			Name:      match[2],
			Content:   content,
			Direction: dir,
		}

		fsm.Migrations[dir] = append(fsm.Migrations[dir], migration)
	}
	return nil

}
