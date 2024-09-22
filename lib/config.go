package lib

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

func IsFirstTimeUser() bool {
	_, err := libUtils.FindConfigFilePath()
	return err != nil
}

func InitConfigFile(v *viper.Viper, configDirPath string, write bool) error {
	configFileNameSplit := strings.Split(constants.SelfConfigFileName, ".")

	v.SetConfigName(strings.Join(configFileNameSplit[:len(configFileNameSplit)-1], "."))
	v.SetConfigType(configFileNameSplit[len(configFileNameSplit)-1])
	v.AddConfigPath(configDirPath)

	connConfig := types.NewDefaultPgConnConfig()
	v.SetDefault(constants.DbConfigConnSubkey, connConfig)
	v.SetDefault(constants.DbConfigCurrentBranchKey, "postgres")
	v.SetDefault(constants.DbConfigOriginsKey, map[string]types.ConnConfig{})

	if write {
		return v.WriteConfigAs(path.Join(configDirPath, constants.SelfConfigFileName))
	}

	return nil
}

func ReadConfig(v *viper.Viper, configFilePath string) error {
	configDirPath, _ := path.Split(configFilePath)
	configFileNameSplit := strings.Split(constants.SelfConfigFileName, ".")

	v.SetConfigName(strings.Join(configFileNameSplit[:len(configFileNameSplit)-1], "."))
	v.SetConfigType(configFileNameSplit[len(configFileNameSplit)-1])
	v.AddConfigPath(configDirPath)

	return v.ReadInConfig()
}

func EnsureSupportedDbDriver() error {
	driver := viper.GetString(constants.DbConfigDriverKey)

	if !slices.Contains(constants.SupportedDrivers, driver) {
		return fmt.Errorf("Unsupported database driver '%s'", driver)
	}

	return nil
}
