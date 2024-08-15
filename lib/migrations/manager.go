package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/libUtils"
	"dbdaddy/types"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
)

func FetchMigrations(dbname string) ([]*types.DbMigration, error) {
	migrations := []*types.DbMigration{}

	migDir, migDirErr := libUtils.GetMigDirPath(dbname)
	if migDirErr != nil {
		return migrations, migDirErr
	}

	dirs, err := os.ReadDir(migDir)
	if err != nil {
		return migrations, err
	}
	slices.SortFunc(dirs, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for i, migDirEntry := range dirs {
		mig := &types.DbMigration{
			DirPath: path.Join(migDir, migDirEntry.Name()),
		}

		files, err := os.ReadDir(migDir)
		if err != nil {
			return migrations, err
		}
		isActive := slices.ContainsFunc(files, func(fileEntry fs.DirEntry) bool {
			return fileEntry.Name() == constants.MigDirActiveSignalFile
		})
		mig.IsActive = isActive

		if i > 0 {
			prevMig := migrations[i-1]
			mig.Down = prevMig
			prevMig.Up = mig
		}

		migrations = append(migrations, mig)
	}

	return migrations, nil
}
