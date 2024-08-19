package sqlwriter

import (
	constants "dbdaddy/const"
	"dbdaddy/sqlwriter/sqlpg"
	"dbdaddy/types"

	"github.com/spf13/viper"
)

func getDriver() string {
	driver := viper.GetString(constants.DbConfigDriverKey)
	return driver
}

func GetDisableConstSQL() string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL DEFERRED;"
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 0;"
	default:
		panic("")
	}
}

func GetEnableConstSQL() string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL IMMEDIATE;"
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 1;"
	default:
		panic("")
	}
}

func GetColDefSQL(col *types.Column) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetColDefSQL(col)
	default:
		panic("unsupported driver")
	}
}

func GetCreateTableSQL(tableSchema *types.TableSchema) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetCreateTableSQL(tableSchema)
	default:
		panic("unsupported driver")
	}
}

// tableid should be of the form "schemaname"."tablename"
func GetDropTableSQL(tableid string) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetDropTableSQL(tableid)
	default:
		panic("unsupported driver")
	}
}

func GetATDropColSQL(tableid string, colName string) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATDropColSQL(tableid, colName)
	default:
		panic("unsupported driver")
	}
}

func GetATCreateColSQL(tableid string, col *types.Column) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATCreateColSQL(tableid, col)
	default:
		panic("unsupported driver")
	}
}

func GetATDropConstraint(tableid string, conName string) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATDropConstraint(tableid, conName)
	default:
		panic("unsupported driver")
	}
}

func GetATCreateConstraintSQL(tableid string, con *types.DbConstraint) string {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATCreateConstraintSQL(tableid, con)
	default:
		panic("unsupported driver")
	}
}
