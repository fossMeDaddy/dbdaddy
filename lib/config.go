package lib

import (
	constants "dbdaddy/const"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitConfigFile(v *viper.Viper) {
	localViper := v
	if localViper == nil {
		localViper = viper.GetViper()
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		configPath = constants.GetGlobalConfigPath()
		EnsureDirExists(constants.GetGlobalDirPath())
	}

	tmp := strings.Split(configPath, "/")
	viperConfigPath := strings.Join(tmp[:len(tmp)-1], "/")

	localViper.SetConfigName("dbdaddy.config")
	localViper.SetConfigType("json")
	localViper.AddConfigPath(viperConfigPath)

	// setup default values for PG connection
	localViper.SetDefault(constants.DbConfigHostKey, "localhost")
	localViper.SetDefault(constants.DbConfigPortKey, 5432)
	localViper.SetDefault(constants.DbConfigUserKey, "postgres")
	localViper.SetDefault(constants.DbConfigPassKey, "postgres")
	localViper.SetDefault(constants.DbConfigDriverKey, constants.DbDriverPostgres)
	localViper.SetDefault(constants.DbConfigDbNameKey, "postgres")
	localViper.SetDefault(constants.DbConfigCurrentBranchKey, "postgres")

	configNotFoundErr := localViper.ReadInConfig()
	if configNotFoundErr != nil {
		// NOT SURE ABOUT THIS ONE
		// FindCwdDBCreds()

		writeErr := localViper.SafeWriteConfig()
		if writeErr != nil {
			fmt.Println("Terribly sorry, this is probably a skill issue!\nI can't generate your dbdaddy.config.json")
			fmt.Println(writeErr)
			os.Exit(1)
		}
	}
}
