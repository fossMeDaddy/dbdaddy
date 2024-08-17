package migrationsLib

import (
	"dbdaddy/libUtils"
	"dbdaddy/types"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
)

// returns: list of migrations, the active migration index (-1 if no active migration found) and "isInit"
// to check if the migrations directory was initialized
func Status(dbname string, currentState *types.DbSchema) (types.MigrationStatus, error) {
	var (
		wg sync.WaitGroup
		mx sync.Mutex
	)

	migStat := types.MigrationStatus{}

	migDir, migDirErr := libUtils.GetMigrationsDir(dbname)
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

	for i, migDirEntry := range dirs {
		mig := types.DbMigration{
			DirPath: path.Join(migDir, migDirEntry.Name()),
		}

		wg.Add(1)
		go (func(mig *types.DbMigration) {
			defer wg.Done()

			state, err := mig.ReadState()
			if err != nil {
				fmt.Println("WARNING: error occured while reading state file in", mig.DirPath)
				fmt.Println(err)
				return
			}

			changes := DiffDbSchema(currentState, state)
			if len(changes) == 0 {
				mig.IsActive = true

				mx.Lock()
				migStat.ActiveMigration = mig
				mx.Unlock()
			}
		})(&mig)

		if i > 0 {
			prevMig := migStat.Migrations[i-1]
			mig.Down = &prevMig
			prevMig.Up = &mig
		}

		migStat.Migrations = append(migStat.Migrations, mig)
	}
	migStat.IsInit = len(migStat.Migrations) == 0

	wg.Wait()

	return migStat, nil
}

// migrations array should be an ascendingly sorted array
// returns -1 if no migration is set active
func GetActiveMigration(migrations []types.DbMigration) int {
	findI := slices.IndexFunc(migrations, func(mig types.DbMigration) bool {
		return mig.IsActive
	})

	return findI
}
