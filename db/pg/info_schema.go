package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
	"dbdaddy/types"

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

func GetTableSchema(schema string, tablename string) (types.TableSchema, error) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	table := types.TableSchema{
		Db:     currBranch,
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
