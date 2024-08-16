package migrationsLib

import (
	constants "dbdaddy/const"
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

func Status(dbname string, currentState *types.DbSchema) ([]types.DbMigration, error) {
	var wg sync.WaitGroup
	migrations := []types.DbMigration{}

	migDir, migDirErr := libUtils.GetMigrationsDir(dbname)
	if migDirErr != nil {
		return migrations, migDirErr
	}

	dirs, err := os.ReadDir(migDir)
	if err != nil {
		return migrations, err
	}
	slices.SortFunc(dirs, func(a, b fs.DirEntry) int {
		// always prefer init file to be "smaller" in compare func
		if a.Name() == constants.MigInitDirName {
			return -1
		}

		if b.Name() == constants.MigInitDirName {
			return 1
		}

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
			}
		})(&mig)

		if i > 0 {
			prevMig := migrations[i-1]
			mig.Down = &prevMig
			prevMig.Up = &mig
		}

		migrations = append(migrations, mig)
	}

	wg.Wait()

	return migrations, nil
}

// migrations array should be an ascendingly sorted array
// returns -1 if no migration is set active
func GetActiveMigration(migrations []types.DbMigration) int {
	findI := slices.IndexFunc(migrations, func(mig types.DbMigration) bool {
		return mig.IsActive
	})

	return findI
}
