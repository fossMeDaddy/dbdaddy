package msqlq

import "fmt"

var QGetTableSchema = func(schema string, table string) string {
	return fmt.Sprintf(`
	SELECT 
		column_name,
		default_value,
		nullable,
		datatype,
		is_primary_key,
		is_foreign_key,
		foreign_table_schema,
		foreign_table_name,
		foreign_column_name
	FROM (
		SELECT 
			info.COLUMN_NAME AS column_name,
			CASE
				WHEN info.COLUMN_DEFAULT IS NULL THEN ''
				ELSE info.COLUMN_DEFAULT
			END AS default_value,
			CASE
				WHEN info.IS_NULLABLE = 'YES' THEN TRUE
				ELSE FALSE
			END AS nullable,
			info.DATA_TYPE AS datatype,
			MAX(CASE
				WHEN rel.CONSTRAINT_NAME = 'PRIMARY' THEN TRUE
				ELSE FALSE
			END) AS is_primary_key,
			MAX(CASE
				WHEN rel.CONSTRAINT_NAME LIKE '%%fk%%' THEN TRUE
				ELSE FALSE
			END) AS is_foreign_key,
			MAX(CASE
				WHEN rel.REFERENCED_TABLE_SCHEMA IS NULL THEN ''
				ELSE rel.REFERENCED_TABLE_SCHEMA
			END) AS foreign_table_schema,
			MAX(CASE
				WHEN rel.REFERENCED_TABLE_NAME IS NULL THEN ''
				ELSE rel.REFERENCED_TABLE_NAME
			END) AS foreign_table_name,
			MAX(CASE
				WHEN rel.REFERENCED_COLUMN_NAME IS NULL THEN ''
				ELSE rel.REFERENCED_COLUMN_NAME
			END) AS foreign_column_name,
			MAX(info.ORDINAL_POSITION) AS ordinal_position
		FROM 
			information_schema.COLUMNS info
		LEFT JOIN 
			information_schema.KEY_COLUMN_USAGE rel 
			ON info.TABLE_SCHEMA = rel.TABLE_SCHEMA 
			AND info.TABLE_NAME = rel.TABLE_NAME 
			AND info.COLUMN_NAME = rel.COLUMN_NAME
		WHERE 
			info.TABLE_SCHEMA = '%s' AND
			info.TABLE_NAME = '%s' 
		GROUP BY 
			info.COLUMN_NAME, info.COLUMN_DEFAULT, info.IS_NULLABLE, info.DATA_TYPE
	) AS subquery
	ORDER BY 
		ORDINAL_POSITION
    `, schema, table)
}
