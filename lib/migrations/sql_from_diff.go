package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/libUtils"
	"dbdaddy/sqlwriter"
	"dbdaddy/types"
	"fmt"
	"slices"
)

func getTableSchemaFromEntity(entity []string, currentState, prevState *types.DbSchema) *types.TableSchema {
	var tableSchema *types.TableSchema

	tableid := libUtils.GetTableId(entity[2], entity[3])
	if entity[0] == currentStateTag {
		tableSchema = currentState.Tables[tableid]
	} else {
		tableSchema = prevState.Tables[tableid]
	}

	return tableSchema
}

func getViewSchemaFromEntity(entity []string, currentState, prevState *types.DbSchema) *types.TableSchema {
	var viewSchema *types.TableSchema

	viewid := libUtils.GetTableId(entity[2], entity[3])
	if entity[0] == currentStateTag {
		viewSchema = currentState.Views[viewid]
	} else {
		viewSchema = prevState.Views[viewid]
	}

	return viewSchema
}

func getConstraintFromEntity(entityId []string, currentState, prevState *types.DbSchema) *types.DbConstraint {
	tableSchema := getTableSchemaFromEntity(entityId, currentState, prevState)
	findI := slices.IndexFunc(tableSchema.Constraints, func(con *types.DbConstraint) bool {
		return con.ConName == entityId[4]
	})
	con := tableSchema.Constraints[findI]

	return con
}

func GetSQLFromDiffChanges(currentState, prevState *types.DbSchema, changes []types.MigAction) string {
	sqlFile := ""
	sqlFile += fmt.Sprintln(sqlwriter.GetDisableConstSQL())
	sqlFile += fmt.Sprintln()

	for _, change := range changes {
		switch change.Type {
		// TABLE CHANGES
		case constants.MigActionDropTable:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			sqlFile += sqlwriter.GetDropTableSQL(libUtils.GetTableId(tableSchema.Schema, tableSchema.Name))
		case constants.MigActionCreateTable:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			sqlFile += sqlwriter.GetCreateTableSQL(tableSchema)

		// TABLE COL CHANGES (ALTER TABLE)
		case constants.MigActionDropCol:
			tableid := libUtils.GetTableId(change.EntityId[2], change.EntityId[3])
			sqlFile += sqlwriter.GetATDropColSQL(tableid, change.EntityId[4])
		case constants.MigActionCreateCol:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			findI := slices.IndexFunc(tableSchema.Columns, func(col types.Column) bool {
				return col.Name == change.EntityId[4]
			})
			sqlFile += sqlwriter.GetATCreateColSQL(
				libUtils.GetTableId(tableSchema.Schema, tableSchema.Name),
				&tableSchema.Columns[findI],
			)

		// CREATE VIEW
		case constants.MigActionDropView:
			viewid := libUtils.GetTableId(change.EntityId[2], change.EntityId[3])
			sqlFile += sqlwriter.GetDropViewSQL(viewid)
		// DROP VIEW
		case constants.MigActionCreateView:
			viewSchema := getViewSchemaFromEntity(change.EntityId, currentState, prevState)
			sqlFile += sqlwriter.GetCreateViewSQL(viewSchema)

		// CONSTRAINT CHANGES (ALTER TABLE)
		case constants.MigActionDropConstraint:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			con := getConstraintFromEntity(change.EntityId, currentState, prevState)
			sqlFile += sqlwriter.GetATDropConstraint(
				libUtils.GetTableId(tableSchema.Schema, tableSchema.Name),
				con.ConName,
			)
		case constants.MigActionCreateConstraint:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			con := getConstraintFromEntity(change.EntityId, currentState, prevState)
			sqlFile += sqlwriter.GetATCreateConstraintSQL(
				libUtils.GetTableId(tableSchema.Schema, tableSchema.Name),
				con,
			)
		}

		sqlFile += fmt.Sprintln()
	}

	sqlFile += sqlwriter.GetEnableConstSQL()
	sqlFile += fmt.Sprintln()

	return sqlFile
}
