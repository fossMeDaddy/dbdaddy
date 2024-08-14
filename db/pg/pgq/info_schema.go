// not necessary to move ALL SQL queries here
package pgq

import "fmt"

var QGetTableSchema = func(schema string, table string) string {
	return fmt.Sprintf(`
        select
            infcol.column_name as column_name,
            CASE
                WHEN infcol.column_default IS NULL then ''
                ELSE infcol.column_default
            END as default_value,
            CASE
                WHEN infcol.is_nullable = 'YES' THEN TRUE
                ELSE FALSE
            END AS nullable,
            CASE
				WHEN infcol.data_type = 'ARRAY'
					THEN CONCAT(releltype.data_type, '[]')
				WHEN infcol.domain_name IS NOT NULL
					THEN infcol.domain_name
				ELSE infcol.udt_name
			END as datatype,
            CASE
				WHEN infcol.character_maximum_length IS NULL THEN -1
				ELSE infcol.character_maximum_length
			END,
			CASE
				WHEN infcol.numeric_precision IS NULL THEN -1
				ELSE infcol.numeric_precision
			END,
			CASE
				WHEN infcol.numeric_scale IS NULL THEN -1
				ELSE infcol.numeric_scale
			END,
            CASE
                WHEN relcon.contype = 'p' then true
                ELSE false
            END as is_primary_key,
            CASE
                WHEN relcon.contype = 'f' then true
                ELSE false
            END as is_relation,
            CASE
                WHEN f_relfcol.table_schema IS NULL then ''
                ELSE f_relfcol.table_schema
            END as foreign_table_schema,
            CASE
                WHEN f_relfcol.table_name IS NULL then ''
                ELSE f_relfcol.table_name
            END as foreign_table_name,
            CASE
                WHEN f_relfcol.column_name IS NULL then ''
                ELSE f_relfcol.column_name
            END as foreign_column_name
        from information_schema.columns infcol

        LEFT JOIN information_schema.element_types releltype
		ON ((infcol.table_catalog, infcol.table_schema, infcol.table_name, 'TABLE', infcol.dtd_identifier)
			= (releltype.object_catalog, releltype.object_schema, releltype.object_name, releltype.object_type, releltype.collection_type_identifier))
        inner join pg_namespace as relnsp on infcol.table_schema = relnsp.nspname
        inner join pg_class as relcls on
            relcls.relname = infcol.table_name and
            relcls.relnamespace = relnsp.oid
        left join pg_constraint as relcon on
            relcon.conrelid = relcls.oid and
            relcon.connamespace = relnsp.oid and
            relcon.conkey[1] = infcol.ordinal_position and
            (relcon.contype = 'p' OR relcon.contype = 'f')


        -- FOREIGN KEYS START --
        left join pg_class as f_relcls on relcon.conrelid = f_relcls.oid
        left join pg_namespace as f_relnsp on f_relcls.relnamespace = f_relnsp.oid

        left join pg_class as f_relfcls on relcon.confrelid = f_relfcls.oid
        left join pg_namespace as f_relfnsp on f_relfcls.relnamespace = f_relfnsp.oid

        left join information_schema.columns as f_relcol on
            f_relcol.table_schema = f_relnsp.nspname and
            f_relcol.table_name = f_relcls.relname and
            f_relcol.ordinal_position = relcon.conkey[1]
        left join information_schema.columns as f_relfcol on
            f_relfcol.table_schema = f_relfnsp.nspname and
            f_relfcol.table_name = f_relfcls.relname and
            f_relfcol.ordinal_position = relcon.confkey[1]
        -- FOREIGN KEYS END --

        where
            infcol.table_schema not in ('pg_catalog', 'information_schema')	and
            infcol.table_schema = '%s' and
            infcol.table_name = '%s'
    `, schema, table)
}

var QGetSchema = func() string {
	return `
        select
            infcol.table_schema,
            infcol.table_name,
            infcol.column_name as column_name,
            CASE
                WHEN infcol.column_default IS NULL then ''
                ELSE infcol.column_default
            END as default_value,
            CASE
                WHEN infcol.is_nullable = 'YES' THEN TRUE
                ELSE FALSE
            END AS nullable,
			CASE
				WHEN infcol.data_type = 'ARRAY'
					THEN CONCAT(releltype.data_type, '[]')
				WHEN infcol.domain_name IS NOT NULL
					THEN infcol.domain_name
				ELSE infcol.udt_name
			END as datatype,
            CASE
				WHEN infcol.character_maximum_length IS NULL THEN -1
				ELSE infcol.character_maximum_length
			END,
			CASE
				WHEN infcol.numeric_precision IS NULL THEN -1
				ELSE infcol.numeric_precision
			END,
			CASE
				WHEN infcol.numeric_scale IS NULL THEN -1
				ELSE infcol.numeric_scale
			END,
            CASE
                WHEN relcon.contype = 'p' then true
                ELSE false
            END as is_primary_key,
            CASE
                WHEN relcon.contype = 'f' then true
                ELSE false
            END as is_relation,
            CASE
                WHEN f_relfcol.table_schema IS NULL then ''
                ELSE f_relfcol.table_schema
            END as foreign_table_schema,
            CASE
                WHEN f_relfcol.table_name IS NULL then ''
                ELSE f_relfcol.table_name
            END as foreign_table_name,
            CASE
                WHEN f_relfcol.column_name IS NULL then ''
                ELSE f_relfcol.column_name
            END as foreign_column_name
        from information_schema.columns infcol

        LEFT JOIN information_schema.element_types releltype
		ON ((infcol.table_catalog, infcol.table_schema, infcol.table_name, 'TABLE', infcol.dtd_identifier)
			= (releltype.object_catalog, releltype.object_schema, releltype.object_name, releltype.object_type, releltype.collection_type_identifier))
		inner join pg_namespace as relnsp on infcol.table_schema = relnsp.nspname
        inner join pg_class as relcls on
            relcls.relname = infcol.table_name and
            relcls.relnamespace = relnsp.oid
        left join pg_constraint as relcon on
            relcon.conrelid = relcls.oid and
            relcon.connamespace = relnsp.oid and
            relcon.conkey[1] = infcol.ordinal_position and
            (relcon.contype = 'p' OR relcon.contype = 'f')


        -- FOREIGN KEYS START --
        left join pg_class as f_relcls on relcon.conrelid = f_relcls.oid
        left join pg_namespace as f_relnsp on relcon.connamespace = f_relnsp.oid

        left join pg_class as f_relfcls on relcon.confrelid = f_relfcls.oid
        left join pg_namespace as f_relfnsp on relcon.connamespace = f_relfnsp.oid

        left join information_schema.columns as f_relcol on
            f_relcol.table_schema = f_relnsp.nspname and
            f_relcol.table_name = f_relcls.relname and
            f_relcol.ordinal_position = relcon.conkey[1]
        left join information_schema.columns as f_relfcol on
            f_relfcol.table_schema = f_relfnsp.nspname and
            f_relfcol.table_name = f_relfcls.relname and
            f_relfcol.ordinal_position = relcon.confkey[1]
        -- FOREIGN KEYS END --

        where
            infcol.table_schema not in ('pg_catalog', 'information_schema')
    `
}
