package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/types"
	"fmt"
	"strings"
)

func GetColDefSQL(col *types.Column) string {
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

	return strings.Join(colSql, " ")
}

func GetCreateTableSQL(tableSchema *types.TableSchema) string {
	sqlStmt := fmt.Sprintf(`CREATE TABLE %s.%s (`, tableSchema.Schema, tableSchema.Name)
	sqlStmt += fmt.Sprintln()

	defer (func() {
		sqlStmt += ");"
	})()

	columns := []string{}

	for _, col := range tableSchema.Columns {
		columns = append(columns, constants.MigFileIndentation+GetColDefSQL(&col))
	}

	sqlStmt += fmt.Sprintln(strings.Join(columns, fmt.Sprintln(",")))
	sqlStmt += fmt.Sprintln(");")
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

// tableid should be of the form "schemaname"."tablename"
func GetDropTableSQL(tableid string) string {
	sqlStmt := fmt.Sprintf("DROP TABLE %s;", tableid)
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

func GetATDropColSQL(tableid string, colName string) string {
	sqlStmt := fmt.Sprintf(`ALTER TABLE %s DROP COLUMN %s;`, tableid, colName)
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

func GetATCreateColSQL(tableid string, col *types.Column) string {
	sqlStmt := fmt.Sprintf(`ALTER TABLE %s`, tableid)
	sqlStmt += fmt.Sprintln()

	sqlStmt += constants.MigFileIndentation
	sqlStmt += fmt.Sprintf("ADD COLUMN %s;", GetColDefSQL(col))
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

func GetATDropConstraint(tableid string, conName string) string {
	sqlStmt := fmt.Sprintf(`ALTER TABLE %s DROP CONSTRAINT %s;`, tableid, conName)
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}

func GetATCreateConstraintSQL(tableid string, con *types.DbConstraint) string {
	sqlStmt := fmt.Sprintf(`ALTER TABLE %s`, tableid)
	sqlStmt += fmt.Sprintln()

	sqlStmt += constants.MigFileIndentation
	sqlStmt += fmt.Sprintf(`ADD CONSTRAINT %s %s;`, con.ConName, con.Syntax)
	sqlStmt += fmt.Sprintln()

	return sqlStmt
}
