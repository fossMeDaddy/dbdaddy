package lib

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"fmt"
	"path"
	"regexp"

	"github.com/spf13/viper"
)

func ValidateBranchName(branchname string) bool {
	re := regexp.MustCompile(`^[\w\-]+$`)
	if re == nil {
		panic("regular expression was not compiled!")
	}

	return re.Match([]byte(branchname))
}

func SetCurrentBranch(branchname string) error {
	if db_int.DbExists(branchname) {
		viper.Set(constants.DbConfigCurrentBranchKey, branchname)
		viper.WriteConfig()
	} else {
		return fmt.Errorf("provided branchname doesn't exist")
	}

	return nil
}

func NewBranchFromCurrent(dbname string) error {
	// dump the db
	configFilePath, _ := FindConfigFilePath()
	dumpFilePath := path.Join(
		GetDriverDumpDir(configFilePath),
		constants.GetDumpFileName(viper.GetString(constants.DbConfigCurrentBranchKey)),
	)
	if err := db_int.DumpDb(dumpFilePath, viper.GetViper()); err != nil {
		return err
	}

	if err := db_int.RestoreDb(dbname, viper.GetViper(), dumpFilePath, true); err != nil {
		return err
	}

	return nil
}
