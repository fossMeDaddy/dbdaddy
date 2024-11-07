package migrationsLib

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/sqlwriter"
	"github.com/fossmedaddy/dbdaddy/types"
)

func GetSQLFromDiffChanges(changes []types.MigAction) (string, error) {
	skipNewLine := false
	sqlFile := ""

	for _, change := range changes {
		switch change.Entity.Type {
		// SCHEMA CHANGES
		case types.EntityTypeSchema:
			schema := change.Entity.Ptr.(*types.Schema)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetDropSchemaSQL(schema)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetCreateSchemaSQL(schema)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// SEQUENCE CHANGES
		case types.EntityTypeSequence:
			seq := change.Entity.Ptr.(*types.DbSequence)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetDropSequenceSQL(seq)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetCreateSequenceSQL(seq)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// TABLE CHANGES
		case types.EntityTypeTable:
			tableSchema := change.Entity.Ptr.(*types.TableSchema)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetDropTableSQL(libUtils.GetTableId(tableSchema.Schema, tableSchema.Name))
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetCreateTableSQL(tableSchema)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// TABLE COL CHANGES (ALTER TABLE)
		case types.EntityTypeColumn:
			col := change.Entity.Ptr.(*types.Column)
			tableid := libUtils.GetTableId(col.TableSchema, col.TableName)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetATDropColSQL(tableid, col.Name)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetATCreateColSQL(tableid, col)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// CREATE/DROP VIEW
		case types.EntityTypeView:
			viewSchema := change.Entity.Ptr.(*types.TableSchema)
			viewid := libUtils.GetTableId(viewSchema.Schema, viewSchema.Name)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetDropViewSQL(viewid)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetCreateViewSQL(viewSchema)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// CONSTRAINT CHANGES (ALTER TABLE)
		case types.EntityTypeConstraint:
			con := change.Entity.Ptr.(*types.DbConstraint)
			tableid := libUtils.GetTableId(con.TableSchema, con.TableName)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetATDropConstraint(tableid, con.ConName)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetATCreateConstraintSQL(tableid, con)
				if err != nil {
					return sqlFile, err
				}
				sqlFile += sqlStr
			}

		// INDEX CHANGES
		case types.EntityTypeIndex:
			ind := change.Entity.Ptr.(*types.DbIndex)
			tableid := libUtils.GetTableId(ind.Schema, ind.TableName)
			if change.ActionType == types.ActionTypeDrop {
				sqlStr, err := sqlwriter.GetDropIndexSQL(tableid, ind.Name)
				if err != nil {
					return sqlFile, err
				}

				sqlFile += sqlStr
			} else {
				sqlStr, err := sqlwriter.GetCreateIndexSQL(ind)
				if err != nil {
					return sqlFile, err
				}

				sqlFile += sqlStr
			}

		default:
			skipNewLine = true
		}

		if !skipNewLine {
			sqlFile += fmt.Sprintln()
		}
	}

	return sqlFile, nil
}
