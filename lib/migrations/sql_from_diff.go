package migrationsLib

import (
	"dbdaddy/libUtils"
	"dbdaddy/sqlwriter"
	"dbdaddy/types"
	"fmt"
)

func GetSQLFromDiffChanges(currentState, prevState *types.DbSchema, changes []types.MigAction) string {
	skipNewLine := false

	sqlFile := ""
	sqlFile += fmt.Sprintln(sqlwriter.GetDisableConstSQL())
	sqlFile += fmt.Sprintln()

	for _, change := range changes {
		switch change.Entity.Type {
		// SCHEMA CHANGES
		case types.EntityTypeSchema:
			schema := change.Entity.Ptr.(*types.Schema)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetDropSchemaSQL(schema)
			} else {
				sqlFile += sqlwriter.GetCreateSchemaSQL(schema)
			}

		// SEQUENCE CHANGES
		case types.EntityTypeSequence:
			seq := change.Entity.Ptr.(*types.DbSequence)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetDropSequenceSQL(seq)
			} else {
				sqlFile += sqlwriter.GetCreateSequenceSQL(seq)
			}

		// TABLE CHANGES
		case types.EntityTypeTable:
			tableSchema := change.Entity.Ptr.(*types.TableSchema)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetDropTableSQL(libUtils.GetTableId(tableSchema.Schema, tableSchema.Name))
			} else {
				sqlFile += sqlwriter.GetCreateTableSQL(tableSchema)
			}

		// TABLE COL CHANGES (ALTER TABLE)
		case types.EntityTypeColumn:
			col := change.Entity.Ptr.(*types.Column)
			tableid := libUtils.GetTableId(col.TableSchema, col.TableName)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetATDropColSQL(tableid, col.Name)
			} else {
				sqlFile += sqlwriter.GetATCreateColSQL(tableid, col)
			}

		// CREATE VIEW
		case types.EntityTypeView:
			viewSchema := change.Entity.Ptr.(*types.TableSchema)
			viewid := libUtils.GetTableId(viewSchema.Schema, viewSchema.Name)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetDropViewSQL(viewid)
			} else {
				sqlFile += sqlwriter.GetCreateViewSQL(viewSchema)
			}

		// CONSTRAINT CHANGES (ALTER TABLE)
		case types.EntityTypeConstraint:
			con := change.Entity.Ptr.(*types.DbConstraint)
			tableid := libUtils.GetTableId(con.TableSchema, con.TableName)
			if change.ActionType == types.ActionTypeDrop {
				sqlFile += sqlwriter.GetATDropConstraint(tableid, con.ConName)
			} else {
				sqlFile += sqlwriter.GetATCreateConstraintSQL(tableid, con)
			}
		default:
			skipNewLine = true
		}

		if !skipNewLine {
			sqlFile += fmt.Sprintln()
		}
	}

	sqlFile += sqlwriter.GetEnableConstSQL()
	sqlFile += fmt.Sprintln()

	return sqlFile
}
