package pgq

import "fmt"

func QGetExistingDbs() string {
	return fmt.Sprintf(
		`select datname from pg_database where datistemplate = false order by datname`,
	)
}

func QCheckDbExists(dbname string) string {
	return fmt.Sprintf(
		`select exists(select 1 from pg_database where datname = '%s')`,
		dbname,
	)
}

func QCreateNewDbFromOldTemplate(newDbName string, oldDbName string) string {
	return fmt.Sprintf(
		`CREATE DATABASE %s TEMPLATE %s`,
		newDbName,
		oldDbName,
	)
}

func QDisconnectAllUsersFromDb(dbname string) string {
	return fmt.Sprintf(
		`
            SELECT pg_terminate_backend(pg_stat_activity.pid)
            FROM pg_stat_activity
            WHERE pg_stat_activity.datname = '%s'
                AND pid <> pg_backend_pid()
        `,
		dbname,
	)
}
