package constants

import (
	"os"
	"path/filepath"
)

func GetGlobalConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfGlobalDirName, "dbdaddy.config.json")
}

func GetLocalConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cwd, "dbdaddy.config.json")
}

func GetGlobalDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfGlobalDirName)
}

const (
	// Self
	SelfGlobalDirName = ".dbdaddy"
	SelfDbName        = "__daddys_home"

	//  Config keys
	DbConfigDriverKey = "connection.driver"
	DbConfigHostKey   = "connection.host"
	DbConfigPortKey   = "connection.port"
	DbConfigDbNameKey = "connection.dbname"
	DbConfigUserKey   = "connection.user"
	DbConfigPassKey   = "connection.password"

	DbConfigCurrentBranchKey = "status.currentBranch"

	// Config possible driver values
	DbDriverPostgres = "postgres"
	DbDriverMySQL    = "mysql"
	DbDriverSqlite   = "sqlite"
)

var SupportedDrivers = []string{DbDriverPostgres}
