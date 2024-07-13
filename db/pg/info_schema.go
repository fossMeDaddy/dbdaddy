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
			infcol.column_name as name,
			CASE
				WHEN infcol.column_default IS NULL then '<null>'
				ELSE infcol.column_default
			END as default_value,
			CASE
				WHEN infcol.is_nullable = 'YES' THEN TRUE
				ELSE FALSE
			END AS nullable,
			infcol.udt_name as datatype,
			CASE
				WHEN relcon.oid IS NULL then false
				ELSE true
			END as is_primary_key
		from information_schema.columns infcol

		inner join pg_namespace as relnsp on infcol.table_schema = relnsp.nspname
		inner join pg_class as relcls on
			relcls.relname = infcol.table_name and
			relcls.relnamespace = relnsp.oid
		left join pg_constraint as relcon on
			relcon.conrelid = relcls.oid and
			relcon.connamespace = relnsp.oid and
			relcon.conkey[1] = infcol.ordinal_position and
			relcon.contype = 'p'

		where
			infcol.table_schema not in ('pg_catalog', 'information_schema') and
			infcol.table_schema = '%s' and
			infcol.table_name = '%s'
	`, schema, tablename))
	if err != nil {
		return table, err
	}

	for rows.Next() {
		column := types.Column{}
		if err := rows.Scan(&column.Name, &column.Default, &column.Nullable, &column.DataType, &column.IsPrimaryKey); err != nil {
			return table, err
		}

		table.Columns = append(table.Columns, column)
	}

	return table, nil
}
