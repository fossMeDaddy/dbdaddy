package db_lib

import (
	constants "dbdaddy/const"
	"dbdaddy/db/pg"
	"fmt"
	"slices"

	"github.com/spf13/viper"
)

// this is a central place for calling driver-specific internal functions

func EnsureSupportedDbDriver() {
	driver := viper.GetString(constants.DbConfigDriverKey)

	if !slices.Contains(constants.SupportedDrivers, driver) {
		panic(fmt.Sprintf("Unsupported database driver '%s'", driver))
	}
}

func GetExistingDbs() ([]string, error) {
	EnsureSupportedDbDriver()

	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.GetExistingDbs()
	} else {
		panic("")
	}
}

func DisconnectAllUsers(dbname string) error {
	EnsureSupportedDbDriver()

	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DisconnectAllUsers(dbname)
	} else {
		panic("")
	}
}

func DbExists(dbname string) bool {
	EnsureSupportedDbDriver()

	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DbExists(dbname)
	} else {
		panic("")
	}
}

func NewDbFromOriginal(originalDbName string, newDbName string) error {
	EnsureSupportedDbDriver()

	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.NewDbFromOriginal(originalDbName, newDbName)
	} else {
		panic("")
	}
}

func DeleteDb(dbname string) error {
	EnsureSupportedDbDriver()

	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DeleteDb(dbname)
	} else {
		panic("")
	}
}
