package db_int

import (
	constants "dbdaddy/const"
	"dbdaddy/db/pg"

	"github.com/spf13/viper"
)

// this is a central place for calling driver-specific internal functions

func GetExistingDbs() ([]string, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.GetExistingDbs()
	} else {
		panic("")
	}
}

func DisconnectAllUsers(dbname string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DisconnectAllUsers(dbname)
	} else {
		panic("")
	}
}

func DbExists(dbname string) bool {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DbExists(dbname)
	} else {
		panic("")
	}
}

func NewDbFromOriginal(originalDbName string, newDbName string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.NewDbFromOriginal(originalDbName, newDbName)
	} else {
		panic("")
	}
}

func CreateDb(dbname string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.CreateDb(dbname)
	} else {
		panic("")
	}
}

func DeleteDb(dbname string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DeleteDb(dbname)
	} else {
		panic("")
	}
}
