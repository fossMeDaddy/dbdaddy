package pg

import (
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
	"dbdaddy/types"
	"fmt"
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

func GetDbSchemaMapping(dbname string) (types.DbSchemaMapping, error) {
	dbSchema := types.DbSchemaMapping{}

	rows, err := db.DB.Query(pgq.QGetSchema())
	if err != nil {
		return dbSchema, err
	}

	tableSchemaMapping := map[string]*types.TableSchema{}

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

	return tableSchemaMapping, nil
}

// runs the query on the current db to get the db schema
func GetDbSchema(dbname string) (types.DbSchema, error) {
	dbSchema := types.DbSchema{}

	tableSchemaMapping, err := GetDbSchemaMapping(dbname)
	if err != nil {
		return dbSchema, err
	}

	for key := range tableSchemaMapping {
		dbSchema = append(dbSchema, *tableSchemaMapping[key])
	}

	return dbSchema, nil
}
