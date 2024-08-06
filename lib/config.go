package lib

import (
	constants "dbdaddy/const"
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

func IsFirstTimeUser() bool {
	_, err := FindConfigFilePath()
	return err != nil
}

func InitConfigFile(v *viper.Viper, configPath string, write bool) {
	tmp := strings.Split(configPath, "/")
	viperConfigPath := strings.Join(tmp[:len(tmp)-1], "/")

	v.SetConfigName("dbdaddy.config")
	v.SetConfigType("json")
	v.AddConfigPath(viperConfigPath)

	// setup default values for PG connection
	v.SetDefault(constants.DbConfigHostKey, "localhost")
	v.SetDefault(constants.DbConfigPortKey, 5432)
	v.SetDefault(constants.DbConfigUserKey, "postgres")
	v.SetDefault(constants.DbConfigPassKey, "postgres")
	v.SetDefault(constants.DbConfigDriverKey, constants.DbDriverPostgres)
	v.SetDefault(constants.DbConfigDbNameKey, "postgres")
	v.SetDefault(constants.DbConfigCurrentBranchKey, "postgres")

	// setup default values for PG connection
	if write {
		v.WriteConfigAs(configPath)
	}
}

func ReadConfig(v *viper.Viper, configPath string) error {
	tmp := strings.Split(configPath, "/")
	viperConfigPath := strings.Join(tmp[:len(tmp)-1], "/")

	v.SetConfigName("dbdaddy.config")
	v.SetConfigType("json")
	v.AddConfigPath(viperConfigPath)

	return v.ReadInConfig()
}

func EnsureSupportedDbDriver() {
	driver := viper.GetString(constants.DbConfigDriverKey)

	if !slices.Contains(constants.SupportedDrivers, driver) {
		panic(fmt.Sprintf("Unsupported database driver '%s'", driver))
	}
}
