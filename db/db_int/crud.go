// this is a central place for calling driver-specific internal functions
package db_int

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/msql"
	"github.com/fossmedaddy/dbdaddy/db/pg"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/types"
)

func getDriver() string {
	return globals.CurrentConnConfig.Driver
}

func GetExistingDbs(showHidden bool) ([]string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.GetExistingDbs(showHidden)
	case constants.DbDriverMySQL:
		return msql.GetExistingDbs(showHidden)
	default:
		return []string{}, errs.ErrUnsupportedDriver
	}
}

func DbExists(dbname string) bool {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.DbExists(dbname)
	case constants.DbDriverMySQL:
		return msql.DbExists(dbname)
	default:
		return false
	}
}

func CreateDb(dbname string) error {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.CreateDb(dbname)
	case constants.DbDriverMySQL:
		return msql.CreateDb(dbname)
	default:
		return errs.ErrUnsupportedDriver
	}
}

func DeleteDb(dbname string) error {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.DeleteDb(dbname)
	case constants.DbDriverMySQL:
		return msql.DeleteDb(dbname)
	default:
		return errs.ErrUnsupportedDriver
	}
}

func ListTablesInDb() ([]types.Table, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.ListTablesInDb()
	case constants.DbDriverMySQL:
		return msql.ListTablesInDb()
	default:
		return []types.Table{}, errs.ErrUnsupportedDriver
	}
}

func GetTableSchema(dbname, schema, tablename string) (*types.TableSchema, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.GetTableSchema(schema, tablename)
	case constants.DbDriverMySQL:
		return msql.GetTableSchema(schema, tablename)
	default:
		return nil, errs.ErrUnsupportedDriver
	}
}

func GetDbSchema() (*types.DbSchema, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return pg.GetDbSchema("", "")
	default:
		return nil, errs.ErrUnsupportedDriver
	}
}
