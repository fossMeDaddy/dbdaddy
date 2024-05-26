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

	return filepath.Join(cwd, SelfGlobalDirName, "dbdaddy.config.json")
}

func GetGlobalDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfGlobalDirName)
}

const (
	DbConfigDefaultPostgresUrl = "postgresql://postgres:postgres@localhost:5342/postgres"

	// Self
	SelfGlobalDirName = ".dbdaddy"

	//  Config keys
	DbConfigDriverKey        = "connection.driver"
	DbConfigHostKey          = "connection.host"
	DbConfigPortKey          = "connection.port"
	DbConfigDbNameKey        = "connection.dbname"
	DbConfigUserKey          = "connection.user"
	DbConfigPassKey          = "connection.password"
	DbConfigCurrentBranchKey = "status.currentBranch"

	// Config possible driver values
	DbDriverPostgres = "postgres"
	DbDriverMySQL    = "mysql"
	DbDriverSqlite   = "sqlite"

	// dumps
	PgDumpDir     = "pg_dumps"
	MySqlDumpDir  = "mysql_dumps"
	SqliteDumpDir = "sqlite_dumps"
)

var DriverDumpDirNames = map[string]string{
	DbDriverPostgres: PgDumpDir,
	DbDriverMySQL:    MySqlDumpDir,
	DbDriverSqlite:   SqliteDumpDir,
}

var SupportedDrivers = []string{DbDriverPostgres}
