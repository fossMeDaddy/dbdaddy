package pg

import (
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
	"dbdaddy/types"
	"fmt"
	"sync"
)

var (
	globalWg sync.WaitGroup
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

func GetTableSchema(dbname string, schema string, tablename string) (types.TableSchema, error) {
	table := types.TableSchema{
		Db:     dbname,
		Schema: schema,
		Name:   tablename,
	}

	rows, err := db.DB.Query(pgq.QGetTableSchema(schema, tablename))
	if err != nil {
		return table, err
	}

	for rows.Next() {
		column := types.Column{}
		if err := rows.Scan(
			&column.Name,
			&column.Default,
			&column.Nullable,
			&column.DataType,
			&column.CharMaxLen,
			&column.IsPrimaryKey,
			&column.IsRelation,
			&column.ForeignTableSchema,
			&column.ForeignTableName,
			&column.ForeignColumnName,
		); err != nil {
			return table, err
		}

		table.Columns = append(table.Columns, column)
	}

	return table, nil
}

func GetDbSchema(dbname string) (types.DbSchema, error) {
	dbSchema := types.DbSchema{}

	rows, err := db.DB.Query(pgq.QGetSchema())
	if err != nil {
		return dbSchema, err
	}

	tableSchemaMapping := map[string]*types.TableSchema{}

	dbTypes := []types.DbType{}
	globalWg.Add(1)
	go (func() {
		defer globalWg.Done()

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
			fmt.Println("Warning, error occured while fetching DB types", err)
			return
		}

		for rows.Next() {
			dbType := types.DbType{}

			err := rows.Scan(
				&dbType.Schema,
				&dbType.Name,
			)
			if err != nil {
				fmt.Println("error occured while reading row from pg_type", err)
				return
			}

			dbTypes = append(dbTypes, dbType)
		}
	})()

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
			&column.IsPrimaryKey,
			&column.IsRelation,
			&column.ForeignTableSchema,
			&column.ForeignTableName,
			&column.ForeignColumnName,
		); err != nil {
			return dbSchema, err
		}

		tableid := fmt.Sprint(tableschema + tablename)

		tableSchema := tableSchemaMapping[tableid]
		if tableSchema == nil {
			tableSchema = &types.TableSchema{
				Db:     dbname,
				Schema: tableschema,
				Name:   tablename,
			}

			tableSchemaMapping[tableid] = tableSchema
		}

		tableSchema.Columns = append(tableSchema.Columns, column)
	}

	globalWg.Wait()

	dbSchema.Types = dbTypes
	dbSchema.Tables = tableSchemaMapping

	return dbSchema, nil
}
