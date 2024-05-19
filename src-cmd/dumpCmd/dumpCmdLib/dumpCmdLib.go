package dumpCmdLib

import (
	"dbdaddy/lib"
	"os"
	"path"
	"slices"
	"strings"
)

type DbDumpFilesMap = map[string][]string

func GetDbGroupedDumpFiles(configPath string) (DbDumpFilesMap, error) {
	dumpDirPath := lib.GetDriverDumpDir(configPath)
	files, err := os.ReadDir(dumpDirPath)
	if err != nil {
		return nil, err
	}

	dumpDbGroups := DbDumpFilesMap{}
	for _, entry := range files {
		fileFullPath := path.Join(dumpDirPath, entry.Name())

		tmp := strings.Split(entry.Name(), "__")
		if len(tmp) < 2 {
			continue
		}

		dbname := tmp[1]
		if dumpDbGroups[dbname] != nil {
			dumpDbGroups[dbname] = append(dumpDbGroups[dbname], fileFullPath)
		} else {
			dumpDbGroups[dbname] = []string{fileFullPath}
		}
	}

	for dbname := range dumpDbGroups {
		slices.Sort(dumpDbGroups[dbname])
		slices.Reverse(dumpDbGroups[dbname])
	}

	return dumpDbGroups, nil
}
