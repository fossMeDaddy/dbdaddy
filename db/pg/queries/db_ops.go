package pg_queries

import "fmt"

// GET_EXISTING_DBS = `select datname from pg_database where datistemplate = false order by datname`

// CHECK_DB_EXISTS = `select exists(select 1 from pg_database where datname = '{{ db_name }}')`

// CREATE_NEW_DB_FROM_OLD_TEMPLATE = `CREATE DATABASE {{ new_db_name }} TEMPLATE {{ db_name }}`

// CREATE_NEW_DB = `CREATE DATABASE {{ db_name }}`

// DISCONNECT_ALL_USERS_FROM_DB = `
// 	SELECT pg_terminate_backend(pg_stat_activity.pid)
// 	FROM pg_stat_activity
// 	WHERE pg_stat_activity.datname = '{{db_name}}'
// 		AND pid <> pg_backend_pid()
// `

// DELETE_DB = `drop database {{ db_name }}`

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

func QCreateNewDb(dbname string) string {
	return fmt.Sprintf(
		`CREATE DATABASE %s`,
		dbname,
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

func QDeleteDb(dbname string) string {
	return fmt.Sprintf(
		`drop database %s`,
		dbname,
	)
}
