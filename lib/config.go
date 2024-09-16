package lib

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"

	"github.com/spf13/viper"
)

func IsFirstTimeUser() bool {
	_, err := libUtils.FindConfigFilePath()
	return err != nil
}

func InitConfigFile(v *viper.Viper, configDirPath string, write bool) {
	configFileNameSplit := strings.Split(constants.SelfConfigFileName, ".")

	v.SetConfigName(strings.Join(configFileNameSplit[:len(configFileNameSplit)-1], "."))
	v.SetConfigType(configFileNameSplit[len(configFileNameSplit)-1])
	v.AddConfigPath(configDirPath)

	// setup default values for PG connection
	v.SetDefault(constants.DbConfigHostKey, "127.0.0.1")
	v.SetDefault(constants.DbConfigPortKey, 5432)
	v.SetDefault(constants.DbConfigUserKey, "postgres")
	v.SetDefault(constants.DbConfigPassKey, "postgres")
	v.SetDefault(constants.DbConfigDriverKey, constants.DbDriverPostgres)
	v.SetDefault(constants.DbConfigDbNameKey, "postgres")
	v.SetDefault(constants.DbConfigCurrentBranchKey, "postgres")

	// setup default values for PG connection
	if write {
		v.WriteConfigAs(path.Join(configDirPath, constants.SelfConfigFileName))
	}
}

func ReadConfig(v *viper.Viper, configFilePath string) error {
	configDirPath, _ := path.Split(configFilePath)
	configFileNameSplit := strings.Split(constants.SelfConfigFileName, ".")

	v.SetConfigName(strings.Join(configFileNameSplit[:len(configFileNameSplit)-1], "."))
	v.SetConfigType(configFileNameSplit[len(configFileNameSplit)-1])
	v.AddConfigPath(configDirPath)

	return v.ReadInConfig()
}

func EnsureSupportedDbDriver() {
	driver := viper.GetString(constants.DbConfigDriverKey)

	if !slices.Contains(constants.SupportedDrivers, driver) {
		panic(fmt.Sprintf("Unsupported database driver '%s'", driver))
	}
}
