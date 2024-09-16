package lib

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"

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

func NewBranchFromCurrent(dbname string, onlySchema bool) error {
	if db_int.DbExists(dbname) {
		return errs.ErrDbAlreadyExists
	}

	configFilePath, _ := libUtils.FindConfigFilePath()
	dumpFilePath := path.Join(
		GetDriverDumpDir(configFilePath),
		constants.GetDumpFileName(viper.GetString(constants.DbConfigCurrentBranchKey)),
	)
	if err := db_int.DumpDb(dumpFilePath, viper.GetViper(), onlySchema); err != nil {
		return err
	}

	if err := db_int.RestoreDb(dbname, viper.GetViper(), dumpFilePath, true); err != nil {
		return err
	}
	if err := os.RemoveAll(dumpFilePath); err != nil {
		return err
	}
	return nil
}
