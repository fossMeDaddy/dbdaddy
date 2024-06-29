package pg

import (
	"dbdaddy/db"
	"dbdaddy/types"
	"fmt"
)

func ListTablesInDb() ([]types.Table, error) {
	tables := []types.Table{}
	rows, err := db.DB.Query(`
		select distinct table_name as name, table_schema as schema from information_schema.tables
        where table_schema != 'information_schema' and table_schema != 'pg_catalog'
        order by table_name
    `)
	if err != nil {
		return tables, err
	}

	for rows.Next() {
		table := types.Table{}
		if err := rows.Scan(&table.Name, &table.Schema); err != nil {
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

	rows, err := db.DB.Query(fmt.Sprintf(`
		select
			column_name as name,
			CASE
				WHEN column_default IS NULL then '<null>'
				ELSE column_default
			END as default_value,
			CASE
				WHEN is_nullable = 'YES' THEN TRUE
				ELSE FALSE
			END AS nullable,
			udt_name as datatype
		from information_schema.columns
			where
				table_schema != 'information_schema' and
				table_schema != 'pg_catalog' and
				table_schema = '%s' and
				table_name = '%s'
	`, schema, tablename))
	if err != nil {
		return table, err
	}

	for rows.Next() {
		column := types.Column{}
		if err := rows.Scan(&column.Name, &column.Default, &column.Nullable, &column.DataType); err != nil {
			return table, err
		}

		table.Columns = append(table.Columns, column)
	}

	return table, nil
}
