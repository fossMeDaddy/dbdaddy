package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/types"
	"fmt"
)

func GetCreateTableSQL(tableSchema *types.TableSchema) string {
	sqlStmt := fmt.Sprintln(fmt.Sprintf(`CREATE TABLE %s.%s (`, tableSchema.Schema, tableSchema.Name))
	defer (func() {
		sqlStmt += ");"
	})()

	for _, col := range tableSchema.Columns {
		// colSql := fmt.Sprintf(`
		//     %s %s
		// `, col.Name, col.DataType, col.))
		col.Name = ""
	}

	return sqlStmt
}

func GetSQLFromDiffChanges(currentState, prevState *types.DbSchema, changes []types.MigAction) string {
	sqlFile := ""
	sqlFile += fmt.Sprintln(db_int.GetDisableConstSQL())

	for _, change := range changes {
		switch change.Type {
		case constants.MigActionCreateTable:
			// sqlFile +=
		}
	}

	return sqlFile
}
