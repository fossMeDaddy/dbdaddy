package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
	"dbdaddy/libUtils"
	"dbdaddy/types"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

func ListTablesInDb() ([]types.Table, error) {
	tables := []types.Table{}
	rows, err := db.DB.Query(`
		select
			table_name as name,
			table_schema as schema,
			table_type as type
		from information_schema.tables

		where
			table_schema not in ('information_schema', 'pg_catalog') and
			table_type in ('BASE TABLE', 'VIEW')

		order by table_name
    `)
	if err != nil {
		return tables, err
	}

	for rows.Next() {
		table := types.Table{}
		if err := rows.Scan(&table.Name, &table.Schema, &table.Type); err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func GetTableSchema(schema string, tablename string) (*types.TableSchema, error) {
	var tableSchema *types.TableSchema

	dbSchema, err := GetDbSchema(schema, tablename)
	if err != nil {
		return tableSchema, err
	}

	tableid := libUtils.GetTableId(schema, tablename)
	tableSchema = dbSchema.Tables[tableid]
	if tableSchema == nil {
		return tableSchema, fmt.Errorf("table with name '%s' could not be found in db schema", tableid)
	}

	return tableSchema, nil
}

// take schema & tablename as "" if need to fetch the whole db schema
// else specify schema & tablename to get schema for a specific table
func GetDbSchema(schema, tablename string) (*types.DbSchema, error) {
	var wg sync.WaitGroup

	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)
	tableid := libUtils.GetTableId(schema, tablename)

	dbSchema := &types.DbSchema{
		DbName: currBranch,
	}

	tableSchemaMapping := map[string]*types.TableSchema{}
	viewSchemaMapping := map[string]*types.TableSchema{}
	schemaPresenceMapping := map[string]bool{}
	dbCons := map[string][]*types.DbConstraint{}
	dbTypes := []types.DbType{}

	wg.Add(2)

	var (
		typeErr error
	)
	go (func() {
		defer wg.Done()

		rows, err := db.DB.Query(`
			select nsp.nspname, typ.typname from pg_type as typ

			inner join pg_namespace as nsp on
				typ.typnamespace = nsp.oid and
				nsp.nspname != 'information_schema' and
				nsp.nspname NOT LIKE 'pg_%'

			where
				typtype not in ('b', 'c')
		`)
		if err != nil {
			typeErr = err
			return
		}

		for rows.Next() {
			dbType := types.DbType{}

			err := rows.Scan(
				&dbType.Schema,
				&dbType.Name,
			)
			if err != nil {
				typeErr = err
				return
			}

			dbTypes = append(dbTypes, dbType)
		}
	})()

	var (
		conErr error
	)
	go (func() {
		defer wg.Done()

		conRows, _conErr := db.DB.Query(pgq.QGetConstraints(tableid))
		if _conErr != nil {
			conErr = _conErr
			return
		}
		for conRows.Next() {
			con := types.DbConstraint{}
			err := conRows.Scan(
				&con.ConName,
				&con.ConSchema,
				&con.Type,
				&con.UpdateActionType,
				&con.DeleteActionType,
				&con.Syntax,
				&con.TableSchema,
				&con.TableName,
				&con.ColName,
				&con.FTableSchema,
				&con.FTableName,
				&con.FColName,
			)
			if err != nil {
				conErr = err
				return
			}

			tableid := libUtils.GetTableId(con.TableSchema, con.TableName)

			if dbCons[tableid] == nil {
				dbCons[tableid] = []*types.DbConstraint{}
			}
			dbCons[tableid] = append(dbCons[tableid], &con)
		}
	})()

	schemaRows, schemaErr := db.DB.Query(pgq.QGetSchema(tableid))
	if schemaErr != nil {
		return dbSchema, schemaErr
	}
	for schemaRows.Next() {
		var tabletype, tableschema, tablename string
		column := types.Column{}
		if err := schemaRows.Scan(
			&tabletype,
			&tableschema,
			&tablename,
			&column.Name,
			&column.Default,
			&column.Nullable,
			&column.DataType,
			&column.CharMaxLen,
			&column.NumericPrecision,
			&column.NumericScale,
		); err != nil {
			return dbSchema, schemaErr
		}

		tableid := libUtils.GetTableId(tableschema, tablename)
		column.TableSchema = tableschema
		column.TableName = tablename
		schemaPresenceMapping[tableschema] = true

		var schemaPtr *map[string]*types.TableSchema
		switch tabletype {
		case constants.TableTypeBaseTable:
			schemaPtr = &tableSchemaMapping
		case constants.TableTypeView:
			schemaPtr = &viewSchemaMapping
		}

		schema := *schemaPtr

		tableSchema := schema[tableid]
		tableCons := dbCons[tableid]
		if tableSchema == nil {
			tableSchema = &types.TableSchema{
				Schema:      tableschema,
				Name:        tablename,
				Constraints: tableCons,
			}

			schema[tableid] = tableSchema
		}

		tableSchema.Columns = append(tableSchema.Columns, column)
	}

	wg.Add(1)
	var (
		viewErr error
	)
	go (func() {
		defer wg.Done()

		rows, err := db.DB.Query(pgq.QGetViews())
		if err != nil {
			viewErr = err
			return
		}
		for rows.Next() {
			var viewschema, viewname, syntax string
			err := rows.Scan(
				&viewschema,
				&viewname,
				&syntax,
			)
			if err != nil {
				viewErr = err
				return
			}

			viewSchema := viewSchemaMapping[libUtils.GetTableId(viewschema, viewname)]
			viewSchema.DefSyntax = syntax
		}
	})()

	wg.Wait()

	if typeErr != nil {
		return dbSchema, typeErr
	}
	if conErr != nil {
		return dbSchema, typeErr
	}
	if schemaErr != nil {
		return dbSchema, schemaErr
	}
	if viewErr != nil {
		return dbSchema, viewErr
	}

	schemaList := []types.Schema{}
	for schemaname := range schemaPresenceMapping {
		schemaList = append(schemaList, types.Schema{
			Name: schemaname,
		})
	}

	dbSchema.Types = dbTypes
	dbSchema.Tables = tableSchemaMapping
	dbSchema.Views = viewSchemaMapping
	dbSchema.Schemas = schemaList

	return dbSchema, nil
}
