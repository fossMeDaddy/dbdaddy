package sqlparser

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/sqlparser/sqlparserpg"
	"github.com/fossmedaddy/dbdaddy/types"
)

func getDriver() string {
	return globals.CurrentConnConfig.Driver
}

// takes in sql string that might contain multiple statements
// outputs the dbSchema from sql string containing create table and constraint queries
func GetDbSchemaFromSQL(sqlStr string) (types.DbSchema, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlparserpg.GetDbSchemaFromSQL(sqlStr)
	default:
		return types.DbSchema{}, errs.ErrUnsupportedDriver
	}
}
