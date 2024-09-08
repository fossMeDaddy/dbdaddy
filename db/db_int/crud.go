// this is a central place for calling driver-specific internal functions
package db_int

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/pg"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

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

func ListTablesInDb() ([]types.Table, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.ListTablesInDb()
	} else {
		panic("")
	}
}

func GetTableSchema(dbname, schema, tablename string) (*types.TableSchema, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.GetTableSchema(schema, tablename)
	} else {
		panic("")
	}
}

func GetDbSchema(dbname string) (*types.DbSchema, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.GetDbSchema("", "")
	} else {
		panic("")
	}
}
