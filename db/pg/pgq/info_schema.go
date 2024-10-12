// not necessary to move ALL SQL queries here
package pgq

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
)

// pass "" in tableid if need to get schema for the whole db
func QGetSchema(tableid string) string {

	whereClause := fmt.Sprintf("infcol.table_schema not in %s", constants.PgExcludedDbsSQLStr)
	if len(tableid) > 0 {
		schemaname, tablename := libUtils.GetTableFromId(tableid)
		whereClause += fmt.Sprintf(`
            and infcol.table_schema = '%s' and infcol.table_name = '%s'
        `, schemaname, tablename)
	}

	return fmt.Sprintf(`
        select
            inftable.table_type,
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
				WHEN
					infcol.numeric_precision IS NOT NULL AND
					infcol.numeric_scale > 0
				THEN infcol.numeric_precision
				ELSE -1
			END as numeric_precision,
			CASE
				WHEN
					infcol.numeric_scale IS NULL OR
					infcol.numeric_scale = 0
				THEN -1
				ELSE infcol.numeric_scale
			END AS numeric_scale
        from information_schema.columns infcol

        LEFT JOIN information_schema.element_types releltype
        ON ((infcol.table_catalog, infcol.table_schema, infcol.table_name, 'TABLE', infcol.dtd_identifier)
            = (releltype.object_catalog, releltype.object_schema, releltype.object_name, releltype.object_type, releltype.collection_type_identifier))
        inner join pg_namespace as relnsp on infcol.table_schema = relnsp.nspname
        inner join pg_class as relcls on
            relcls.relname = infcol.table_name and
            relcls.relnamespace = relnsp.oid
	left join information_schema.tables as inftable on
		inftable.table_schema = infcol.table_schema and
		inftable.table_name = infcol.table_name

        where
            %s
	
	order by ordinal_position
    `, whereClause)
}

func QGetViews() string {
	return `
        select
            table_schema,
            table_name,
            view_definition
        from information_schema.views
        where
          table_schema not in ('pg_catalog', 'information_schema') and
          view_definition IS NOT NULL
    `
}

func QGetSequences() string {
	return `
		select
			infseq.sequence_schema,
			infseq.sequence_name,
			infseq.data_type,
			pgseq.increment_by,
			pgseq.min_value,
			pgseq.max_value,
			pgseq.start_value,
			pgseq.cache_size,
			CASE
				WHEN infseq.cycle_option = 'NO' then false
				ELSE true
			END as cycle_option
		from information_schema.sequences as infseq
		inner join pg_sequences as pgseq on infseq.sequence_name = pgseq.sequencename
		where sequence_schema not in ('pg_catalog', 'information_schema')
	`
}
