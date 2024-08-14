package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/types"
	"fmt"
	"strings"
)

func GetCreateTableSQL(tableSchema *types.TableSchema) string {
	sqlStmt := fmt.Sprintln(fmt.Sprintf(`CREATE TABLE %s.%s (`, tableSchema.Schema, tableSchema.Name))
	defer (func() {
		sqlStmt += ");"
	})()

	columns := []string{}

	for _, col := range tableSchema.Columns {
		// name
		colSql := []string{col.Name}

		// type
		if col.CharMaxLen != -1 {
			colSql = append(colSql, fmt.Sprintf("%s(%d)", col.DataType, col.CharMaxLen))
		} else if col.NumericPrecision != -1 {
			bracketContent := fmt.Sprintf("%d", col.NumericPrecision)
			if col.NumericScale != -1 {
				bracketContent += fmt.Sprintf(",%d", col.NumericScale)
			}

			colSql = append(colSql, fmt.Sprintf("%s(%s)", col.DataType, bracketContent))
		} else {
			colSql = append(colSql, col.DataType)
		}

		// default
		if len(col.Default) > 0 {
			colSql = append(colSql, fmt.Sprintf("DEFAULT %s", col.Default))
		}

		// nullable
		if !col.Nullable {
			colSql = append(colSql, "NOT NULL")
		}

		columns = append(columns, "    "+strings.Join(colSql, " "))
	}

	sqlStmt += fmt.Sprintln(strings.Join(columns, fmt.Sprintln(",")))
	sqlStmt += fmt.Sprintln(");")
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

func GetSQLFromDiffChanges(currentState, prevState *types.DbSchema, changes []types.MigAction) string {
	sqlFile := ""
	sqlFile += fmt.Sprintln(db_int.GetDisableConstSQL())
	sqlFile += fmt.Sprintln()

	for _, change := range changes {
		switch change.Type {
		case constants.MigActionCreateTable:
			var tableSchema *types.TableSchema
			if change.EntityId[0] == currentStateTag {
				tableid := strings.Join(change.EntityId[2:], ".")
				tableSchema = currentState.Tables[tableid]
			} else {
				tableid := strings.Join(change.EntityId[2:], ".")
				tableSchema = prevState.Tables[tableid]
			}

			sqlFile += GetCreateTableSQL(tableSchema)
		}
	}

	return sqlFile
}
