package pg

import (
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
	"dbdaddy/libUtils"
	"dbdaddy/types"
	"fmt"
	"sync"
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

func GetTableSchema(dbname string, schema string, tablename string) (*types.TableSchema, error) {
	var tableSchema *types.TableSchema

	dbSchema, err := GetDbSchema(dbname, schema, tablename)
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

func GetDbSchema(dbname, schema, tablename string) (*types.DbSchema, error) {
	var wg sync.WaitGroup

	tableid := libUtils.GetTableId(schema, tablename)

	dbSchema := &types.DbSchema{}

	tableSchemaMapping := map[string]*types.TableSchema{}
	dbCons := map[string][]*types.DbConstraint{}
	dbTypes := []types.DbType{}

	wg.Add(1)
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

	conRows, conErr := db.DB.Query(pgq.QGetConstraints(tableid))
	if conErr != nil {
		return dbSchema, conErr
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
			fmt.Println(err)
			return dbSchema, err
		}

		tableid := libUtils.GetTableId(con.TableSchema, con.TableName)

		if dbCons[tableid] == nil {
			dbCons[tableid] = []*types.DbConstraint{}
		}
		dbCons[tableid] = append(dbCons[tableid], &con)
	}

	rows, err := db.DB.Query(pgq.QGetSchema(tableid))
	if err != nil {
		return dbSchema, err
	}

	for rows.Next() {
		var tableschema, tablename string
		column := types.Column{}
		if err := rows.Scan(
			&tableschema,
			&tablename,
			&column.Name,
			&column.Default,
			&column.Nullable,
			&column.DataType,
			&column.CharMaxLen,
			&column.NumericPrecision,
			&column.NumericScale,
			&column.IsPrimaryKey,
		); err != nil {
			return dbSchema, err
		}

		tableid := libUtils.GetTableId(tableschema, tablename)

		tableSchema := tableSchemaMapping[tableid]
		tableCons := dbCons[tableid]
		if tableSchema == nil {
			tableSchema = &types.TableSchema{
				Db:          dbname,
				Schema:      tableschema,
				Name:        tablename,
				Constraints: tableCons,
			}

			tableSchemaMapping[tableid] = tableSchema
		}

		tableSchema.Columns = append(tableSchema.Columns, column)
	}

	wg.Wait()

	if typeErr != nil {
		return dbSchema, typeErr
	}
	if conErr != nil {
		return dbSchema, typeErr
	}

	dbSchema.Types = dbTypes
	dbSchema.Tables = tableSchemaMapping

	return dbSchema, nil
}
