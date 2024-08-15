package libUtils

import (
	constants "dbdaddy/const"
	"path"
)

func GetMigDirPath(dbname string) (string, error) {
	var migDirPath string

	dirPath, _ := FindConfigDirPath()

	configParentDir := path.Dir(dirPath)
	migDirPath = path.Join(configParentDir, constants.MigDirName, dbname)
	_, err := EnsureDirExists(migDirPath)
	if err != nil {
		return migDirPath, err
	}

	return migDirPath, nil
}
