package libUtils

import (
	constants "dbdaddy/const"
	"path"
	"strings"
	"time"
)

func GetMigrationsDir(dbname string) (string, error) {
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

// returns "<timestamp>__title" or "<timestamp>"
func GenerateMigrationId(title string) string {
	ts := time.Now().Format("2006-01-01T03_04_05")
	fileTitle := strings.ReplaceAll(strings.Trim(strings.ToLower(title), " "), " ", "_")

	if len(fileTitle) > 0 {
		return ts + "__" + fileTitle
	}

	return ts
}
