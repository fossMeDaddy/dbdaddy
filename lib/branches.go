package lib

import (
	"fmt"
	"regexp"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"

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
