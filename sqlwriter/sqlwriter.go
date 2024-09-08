package sqlwriter

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/sqlwriter/sqlpg"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

func getDriver() string {
	driver := viper.GetString(constants.DbConfigDriverKey)
	return driver
}

func GetDisableConstSQL() (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL DEFERRED;", nil
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 0;", nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetEnableConstSQL() (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL IMMEDIATE;", nil
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 1;", nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetDropSchemaSQL(schema *types.Schema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetDropSchemaSQL(schema), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetCreateSchemaSQL(schema *types.Schema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetCreateSchemaSQL(schema), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetDropSequenceSQL(seq *types.DbSequence) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetDropSequenceSQL(seq), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetCreateSequenceSQL(seq *types.DbSequence) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetCreateSequenceSQL(seq), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetColDefSQL(col *types.Column) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetColDefSQL(col), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetCreateTableSQL(tableSchema *types.TableSchema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetCreateTableSQL(tableSchema), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

// tableid should be of the form "schemaname"."tablename"
func GetDropTableSQL(tableid string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetDropTableSQL(tableid), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetCreateViewSQL(viewSchema *types.TableSchema) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetCreateViewSQL(viewSchema), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetDropViewSQL(viewid string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetDropViewSQL(viewid), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetATDropColSQL(tableid string, colName string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATDropColSQL(tableid, colName), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetATCreateColSQL(tableid string, col *types.Column) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATCreateColSQL(tableid, col), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetATDropConstraint(tableid string, conName string) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATDropConstraint(tableid, conName), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}

func GetATCreateConstraintSQL(tableid string, con *types.DbConstraint) (string, error) {
	switch getDriver() {
	case constants.DbDriverPostgres:
		return sqlpg.GetATCreateConstraintSQL(tableid, con), nil
	default:
		return "", fmt.Errorf("unsupported driver")
	}
}
