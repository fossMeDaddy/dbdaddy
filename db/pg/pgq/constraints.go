package pgq

import (
	constants "dbdaddy/const"
	"fmt"
)

// pass "" in tableid if need to get constraints for the whole db
func QGetAllConstraints(tableid string) string {

	whereClause := "infcol.table_schema not in ('pg_catalog', 'information_schema')"
	if len(tableid) > 0 {
		schemaname, tablename := constants.GetTableFromId(tableid)
		whereClause += fmt.Sprintf(`
            and infcol.table_schema = '%s' and infcol.table_name = '%s'
        `, schemaname, tablename)
	}

	return fmt.Sprintf(`
        select
			relcon.conname,
			relconnsp.nspname as connamespace,
			relcon.contype,
			relcon.confupdtype,
			relcon.confdeltype,
			CASE
				WHEN relcon.conbin IS NOT NULL THEN pg_get_constraintdef(relcon.oid, true)
				ELSE ''
			END as check_syntax,
			infcol.table_schema,
            infcol.table_name,
            infcol.column_name as column_name,
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

		inner join pg_namespace as relnsp on infcol.table_schema = relnsp.nspname
        inner join pg_class as relcls on
            relcls.relname = infcol.table_name and
            relcls.relnamespace = relnsp.oid
        inner join pg_constraint as relcon on
            relcon.conrelid = relcls.oid and
            relcon.connamespace = relnsp.oid and
            relcon.conkey[1] = infcol.ordinal_position and
            relcon.contype in ('p', 'f', 'c', 'u')
		inner join pg_namespace as relconnsp on relcon.connamespace = relconnsp.oid
		
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
            %s
    `, whereClause)
}
