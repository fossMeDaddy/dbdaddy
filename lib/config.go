package lib

import (
	constants "dbdaddy/const"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

func OpenConfigFileAt(configFilePath string) {
	vimOsCmd := exec.Command("vim", configFilePath)
	vimOsCmd.Stdin = os.Stdin
	vimOsCmd.Stdout = os.Stdout
	vimOsCmd.Stderr = os.Stderr

	vimErr := vimOsCmd.Run()
	if vimErr != nil {
		fmt.Println("Failed to open vim, trying nano...")

		nanoOsCmd := exec.Command("nano", configFilePath)
		nanoOsCmd.Stdin = os.Stdin
		nanoOsCmd.Stdout = os.Stdout
		nanoOsCmd.Stderr = os.Stderr

		nanoErr := nanoOsCmd.Run()
		if nanoErr != nil {
			fmt.Println("Holy shit bro?! wtf are you using for an OS? no vim, no nano, is this what, an alpine docker container?")
			fmt.Println("nano command gave the error:\n" + nanoErr.Error())
			fmt.Println("vim command gave the error:\n" + vimErr.Error())
			os.Exit(1)
		}
	}
}

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
