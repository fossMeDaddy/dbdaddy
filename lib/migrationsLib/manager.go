package migrationsLib

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/types"
)

// returns: list of migrations, the active migration index (-1 if no active migration found) and "isInit"
// to check if the migrations directory was initialized
func Status(dbname string, currentState *types.DbSchema) (MigrationStatus, error) {
	migStat := MigrationStatus{}

	configFilePath, _ := libUtils.FindConfigFilePath()

	migDir := libUtils.GetMigrationsDir(path.Dir(configFilePath), dbname)
	_, migDirErr := libUtils.EnsureDirExists(migDir)
	if migDirErr != nil {
		return migStat, migDirErr
	}

	dirs, err := os.ReadDir(migDir)
	if err != nil {
		return migStat, err
	}
	slices.SortFunc(dirs, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	activeI := -1
	for i, migDirEntry := range dirs {
		mig := DbMigration{
			DirPath: path.Join(migDir, migDirEntry.Name()),
		}

		state, err := mig.ReadState()
		if err != nil {
			fmt.Println("WARNING: error occured while reading state file in", mig.DirPath)
			fmt.Println(err)
			return migStat, err
		}

		changes := DiffDBSchema(currentState, state)
		if len(changes) == 0 {
			mig.IsActive = true

			activeI = i
		}

		if i > 0 {
			prevMig := &migStat.Migrations[i-1]
			mig.Down = prevMig
			prevMig.Up = &mig
		}

		migStat.Migrations = append(migStat.Migrations, mig)
	}

	if activeI != -1 {
		migStat.ActiveMigration = &migStat.Migrations[activeI]
	}

	migStat.IsInit = len(migStat.Migrations) == 0

	return migStat, nil
}

func ApplyMigrationSQL(migStat MigrationStatus, isUpMigration bool) error {
	if migStat.ActiveMigration == nil {
		return fmt.Errorf("no active migration found, please run 'migrations generate' to update migrations")
	}

	var sqlStr string
	if isUpMigration {
		if migStat.ActiveMigration.Up == nil {
			return fmt.Errorf("already on the latest migration, can't go into future")
		}

		upSql, err := migStat.ActiveMigration.GetUpQuery()
		if err != nil {
			return err
		}

		sqlStr = upSql
	} else {
		if migStat.ActiveMigration.Down == nil {
			return fmt.Errorf("can't go down from here, on the initial migration")
		}

		downSql, err := migStat.ActiveMigration.GetDownQuery()
		if err != nil {
			return err
		}

		sqlStr = downSql
	}

	stmts := libUtils.GetSQLStmts(sqlStr)
	if err := db_int.ExecuteStatementsTx(stmts); err != nil {
		return err
	}

	return nil
}
