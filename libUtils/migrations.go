package libUtils

import (
	constants "dbdaddy/const"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
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

// returns "00069__title" or "00069"
func GenerateMigrationId(migrationsDirPath string, title string) (string, error) {
	dirs, err := os.ReadDir(migrationsDirPath)
	if err != nil {
		return "", err
	}

	latestN := 0
	for _, dir := range dirs {
		n, err := strconv.Atoi(dir.Name())
		if err != nil {
			return "", err
		}

		latestN = max(latestN, n+1)
	}

	fileNum := fmt.Sprintf("%05d", latestN)

	fileTitle := strings.ReplaceAll(strings.Trim(strings.ToLower(title), " "), " ", "_")

	if len(fileTitle) > 0 {
		return fileNum + "__" + fileTitle, nil
	}

	return fileNum, nil
}
