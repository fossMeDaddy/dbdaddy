package types

import (
	constants "dbdaddy/const"
	"dbdaddy/libUtils"
	"encoding/json"
	"os"
	"path"
)

type DbMigration struct {
	// dir path will also contain the id (cwd/migrations/dbname/2024_01_01:00:00:00)
	DirPath  string
	IsActive bool
	Up       *DbMigration
	Down     *DbMigration
}

func (mig *DbMigration) ReadState() (*DbSchema, error) {
	var state *DbSchema

	stateFileB, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirStateFile))
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(stateFileB, &state); err != nil {
		return state, err
	}

	return state, nil
}

func (mig *DbMigration) GetUpQuery() (string, error) {
	file, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirUpSqlFile))
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func (mig *DbMigration) GetDownQuery() (string, error) {
	file, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirDownSqlFile))
	if err != nil {
		return "", nil
	}

	return string(file), nil
}

// set "this" migration as active.
// removes the active signal file from all other migrations in the parent directory
// then creates the active signal file here
func (mig *DbMigration) SetActive() error {
	migrationsDir := path.Dir(mig.DirPath)

	dirs, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		activeSignalFile := path.Join(dir.Name(), constants.MigDirActiveSignalFile)
		if libUtils.Exists(activeSignalFile) {
			if err := os.Remove(activeSignalFile); err != nil {
				return err
			}
		}
	}

	if err := os.WriteFile(path.Join(mig.DirPath, constants.MigDirActiveSignalFile), []byte{}, 0644); err != nil {
		return err
	}

	return nil
}
