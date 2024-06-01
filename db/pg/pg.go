package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"

	"github.com/hoisie/mustache"
)

var (
	GET_EXISTING_DBS = `select datname from pg_database where datistemplate = false order by datname`

	CHECK_DB_EXISTS = `select exists(select 1 from pg_database where datname = '{{ db_name }}')`

	CREATE_NEW_DB_FROM_OLD_TEMPLATE = `CREATE DATABASE {{ new_db_name }} TEMPLATE {{ db_name }}`

	CREATE_NEW_DB = `CREATE DATABASE {{ db_name }}`

	DISCONNECT_ALL_USERS_FROM_DB = `
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '{{db_name}}'
			AND pid <> pg_backend_pid()
	`

	DELETE_DB = `drop database {{ db_name }}`
)

func GetExistingDbs() ([]string, error) {
	getExistingDbsQueryStr := mustache.Render(GET_EXISTING_DBS)
	rows, err := db.DB.Query(getExistingDbsQueryStr)
	if err != nil {
		return nil, err
	}

	existingDbs := []string{}
	for rows.Next() {
		existingDb := ""
		_ = rows.Scan(&existingDb)
		if existingDb == constants.SelfDbName {
			continue
		}

		existingDbs = append(existingDbs, existingDb)
	}

	return existingDbs, nil
}

func DisconnectAllUsers(dbname string) error {
	disconnectQueryStr := mustache.Render(DISCONNECT_ALL_USERS_FROM_DB, map[string]string{"db_name": dbname})

	_, err := db.DB.Query(disconnectQueryStr)
	return err
}

func DbExists(dbname string) bool {
	checkDbExistsQueryStr := mustache.Render(CHECK_DB_EXISTS, map[string]string{"db_name": dbname})
	row := db.DB.QueryRow(checkDbExistsQueryStr)

	exists := false
	_ = row.Scan(&exists)

	return exists
}

func NewDbFromOriginal(originalDbName string, newDbName string) error {
	createQueryStr := mustache.Render(CREATE_NEW_DB_FROM_OLD_TEMPLATE, map[string]string{"new_db_name": newDbName, "db_name": originalDbName})

	err := DisconnectAllUsers(originalDbName)
	if err != nil {
		return err
	}

	_, err = db.DB.Query(createQueryStr)
	if err != nil {
		return err
	}

	return nil
}

func CreateDb(dbname string) error {
	createQueryStr := mustache.Render(CREATE_NEW_DB, map[string]string{"db_name": dbname})
	_, err := db.DB.Query(createQueryStr)
	return err
}

func DeleteDb(dbname string) error {
	deleteDbQueryStr := mustache.Render(DELETE_DB, map[string]string{"db_name": dbname})

	if err := DisconnectAllUsers(dbname); err != nil {
		return err
	}

	_, err := db.DB.Query(deleteDbQueryStr)
	if err != nil {
		return err
	}

	return nil
}
