package migrationsLib

import (
	"encoding/json"
	"os"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/types"
)

type MigrationStatus struct {
	Migrations      []DbMigration
	ActiveMigration *DbMigration
	IsInit          bool
}

/*
primitive DbMigration type, only containing helper function, use migrationsLib for complete use

a db migration directory will look like: /migrations/dbname/2024_01_01T00_00_00

this directory will contain files:
"active" (present ONLY IF migration is currently active),
"state.json" (a state snapshot of the CURRENT state)
"up.sql" (sql file to go up to the next state) [COULD BE ABSENT IF LATEST STATE]
"down.sql" (sql file to go down to the previous state) [ABSENT ON INIT MIGRATION]
"info.md" (user maintained description of the current state changes relative to previous state)
*/
type DbMigration struct {
	// dir path will also contain the id (.../migrations/dbname/<timestamp>)
	DirPath  string
	IsActive bool
	Up       *DbMigration
	Down     *DbMigration
}

func (mig *DbMigration) ReadState() (*types.DbSchema, error) {
	var state *types.DbSchema

	stateFileB, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirStateFile))
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(stateFileB, &state); err != nil {
		return state, err
	}

	return state, nil
}

func (mig *DbMigration) writeSqlFile(sqlFilePath, sql string) error {
	if len(sql) == 0 {
		os.Remove(sqlFilePath)
		return nil
	}

	return os.WriteFile(sqlFilePath, []byte(sql), 0644)
}

func (mig *DbMigration) GetUpQuery() (string, error) {
	file, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirUpSqlFile))
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func (mig *DbMigration) WriteUpQuery(sql string) error {
	upSqlFilePath := path.Join(mig.DirPath, constants.MigDirUpSqlFile)
	return mig.writeSqlFile(upSqlFilePath, sql)
}

func (mig *DbMigration) GetDownQuery() (string, error) {
	file, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirDownSqlFile))
	if err != nil {
		return "", nil
	}

	return string(file), nil
}

func (mig *DbMigration) WriteDownQuery(sql string) error {
	downSqlFilePath := path.Join(mig.DirPath, constants.MigDirDownSqlFile)
	return mig.writeSqlFile(downSqlFilePath, sql)
}

func (mig *DbMigration) GetInfoFile() (string, error) {
	infoFileB, err := os.ReadFile(path.Join(mig.DirPath, constants.MigDirInfoFile))
	if err != nil {
		return "", err
	}

	return string(infoFileB), nil
}

func (mig *DbMigration) WriteInfoFile(infoStr string) error {
	return os.WriteFile(path.Join(mig.DirPath, constants.MigDirInfoFile), []byte(infoStr), 0644)
}

// DEPRECATED
// set "this" migration as active.
// removes the active signal file from all other migrations in the parent directory
// then creates the active signal file in the current migration
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

// ensures migration dir exists and writes state, up/down sql queries & info file to disk
// path string should look like this ".../migrations/dbname/2024-01-01T00_00_00"
func NewDbMigration(migDirPath string, state *types.DbSchema, upSqlStr string, downSqlStr string, infoStr string) (*DbMigration, error) {
	mig := &DbMigration{
		DirPath: migDirPath,
	}

	if _, err := libUtils.EnsureDirExists(migDirPath); err != nil {
		return mig, err
	}

	// write to state file
	stateB, err := json.Marshal(state)
	if err != nil {
		return mig, err
	}
	if err := os.WriteFile(path.Join(migDirPath, constants.MigDirStateFile), stateB, 0644); err != nil {
		return mig, err
	}

	// write to up & down query files
	if err := mig.WriteUpQuery(upSqlStr); err != nil {
		return mig, err
	}

	if err := mig.WriteDownQuery(downSqlStr); err != nil {
		return mig, err
	}

	// write info file blank
	if err := mig.WriteInfoFile(infoStr); err != nil {
		return mig, err
	}

	return mig, nil
}
