package sqlwriter

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/sqlwriter/sqlwriterpg"
	"github.com/fossmedaddy/dbdaddy/types"
)

func getDriver() string {
	return globals.CurrentConnConfig.Driver
}

func GetDisableConstSQL() (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL DEFERRED;", nil
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 0;", nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetEnableConstSQL() (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL IMMEDIATE;", nil
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 1;", nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetDropSchemaSQL(schema *types.Schema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetDropSchemaSQL(schema), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetCreateSchemaSQL(schema *types.Schema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetCreateSchemaSQL(schema), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetDropSequenceSQL(seq *types.DbSequence) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetDropSequenceSQL(seq), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetCreateSequenceSQL(seq *types.DbSequence) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetCreateSequenceSQL(seq), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetColDefSQL(col *types.Column) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetColDefSQL(col), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetCreateTableSQL(tableSchema *types.TableSchema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetCreateTableSQL(tableSchema), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

// tableid should be of the form "schemaname"."tablename"
func GetDropTableSQL(tableid string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetDropTableSQL(tableid), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetCreateViewSQL(viewSchema *types.TableSchema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetCreateViewSQL(viewSchema), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetDropViewSQL(viewid string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetDropViewSQL(viewid), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetATDropColSQL(tableid string, colName string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetATDropColSQL(tableid, colName), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetATCreateColSQL(tableid string, col *types.Column) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetATCreateColSQL(tableid, col), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetATDropConstraint(tableid string, conName string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetATDropConstraint(tableid, conName), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}

func GetATCreateConstraintSQL(tableid string, con *types.DbConstraint) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlwriterpg.GetATCreateConstraintSQL(tableid, con), nil
	default:
		return "", errs.ErrUnsupportedDriver
	}
}
