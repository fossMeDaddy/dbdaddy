// this is a central place for calling driver-specific internal functions
package db_int

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/msql"
	"github.com/fossmedaddy/dbdaddy/db/pg"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

func GetExistingDbs() ([]string, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.GetExistingDbs()
	case constants.DbDriverMySQL:
		return msql.GetExistingDbs()
	default:
		return []string{}, errs.ErrUnsupportedDriver
	}
}

func DbExists(dbname string) bool {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.DbExists(dbname)
	case constants.DbDriverMySQL:
		return msql.DbExists(dbname)
	default:
		return false
	}
}

func CreateDb(dbname string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.CreateDb(dbname)
	case constants.DbDriverMySQL:
		return msql.CreateDb(dbname)
	default:
		return errs.ErrUnsupportedDriver
	}
}

func DeleteDb(dbname string) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.DeleteDb(dbname)
	case constants.DbDriverMySQL:
		return msql.DeleteDb(dbname)
	default:
		return errs.ErrUnsupportedDriver
	}
}

func ListTablesInDb() ([]types.Table, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.ListTablesInDb()
	case constants.DbDriverMySQL:
		return msql.ListTablesInDb()
	default:
		return []types.Table{}, errs.ErrUnsupportedDriver
	}
}

func GetTableSchema(dbname, schema, tablename string) (*types.TableSchema, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.GetTableSchema(schema, tablename)
	case constants.DbDriverMySQL:
		return msql.GetTableSchema(schema, tablename)
	default:
		return nil, errs.ErrUnsupportedDriver
	}
}

func GetDbSchema(dbname string) (*types.DbSchema, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.GetDbSchema("", "")
	default:
		return nil, errs.ErrUnsupportedDriver
	}
}
