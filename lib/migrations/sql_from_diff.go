package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/types"
	"fmt"
	"slices"
)

func getTableSchemaFromEntity(entity []string, currentState, prevState *types.DbSchema) *types.TableSchema {
	var tableSchema *types.TableSchema

	tableid := constants.GetTableId(entity[2], entity[3])
	if entity[0] == currentStateTag {
		tableSchema = currentState.Tables[tableid]
	} else {
		tableSchema = prevState.Tables[tableid]
	}

	return tableSchema
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
	sqlFile += fmt.Sprintln(db_int.GetDisableConstSQL())
	sqlFile += fmt.Sprintln()

	for _, change := range changes {
		switch change.Type {
		// TABLE CHANGES
		case constants.MigActionDropTable:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			sqlFile += GetDropTableSQL(constants.GetTableId(tableSchema.Schema, tableSchema.Name))
		case constants.MigActionCreateTable:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			sqlFile += GetCreateTableSQL(tableSchema)

		// TABLE COL CHANGES (ALTER TABLE)
		case constants.MigActionDropCol:
			tableid := constants.GetTableId(change.EntityId[2], change.EntityId[3])
			sqlFile += GetATDropColSQL(tableid, change.EntityId[4])
		case constants.MigActionCreateCol:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			findI := slices.IndexFunc(tableSchema.Columns, func(col types.Column) bool {
				return col.Name == change.EntityId[4]
			})
			sqlFile += GetATCreateColSQL(
				constants.GetTableId(tableSchema.Schema, tableSchema.Name),
				&tableSchema.Columns[findI],
			)

		// CONSTRAINT CHANGES (ALTER TABLE)
		case constants.MigActionDropConstraint:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			con := getConstraintFromEntity(change.EntityId, currentState, prevState)
			sqlFile += GetATDropConstraint(
				constants.GetTableId(tableSchema.Schema, tableSchema.Name),
				con.ConName,
			)
		case constants.MigActionCreateConstraint:
			tableSchema := getTableSchemaFromEntity(change.EntityId, currentState, prevState)
			con := getConstraintFromEntity(change.EntityId, currentState, prevState)
			sqlFile += GetATCreateConstraintSQL(
				constants.GetTableId(tableSchema.Schema, tableSchema.Name),
				con,
			)
		}

		sqlFile += fmt.Sprintln()
	}

	sqlFile += db_int.GetEnableConstSQL()
	sqlFile += fmt.Sprintln()

	return sqlFile
}
