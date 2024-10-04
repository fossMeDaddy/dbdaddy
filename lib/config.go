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
	v.SetDefault(constants.DbConfigConnKey, connConfig)
	v.SetDefault(constants.DbConfigCurrentBranchKey, connConfig.Database)
	v.SetDefault(constants.DbConfigOriginsKey, types.DbConfigOrigins{})

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
	connConfig, err := libUtils.GetConnConfigFromViper(viper.GetViper())
	if err != nil {
		return err
	}

	if !slices.Contains(constants.SupportedDrivers, connConfig.Driver) {
		return fmt.Errorf("Unsupported database driver '%s'", connConfig.Driver)
	}

	return nil
}
