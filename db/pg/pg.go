package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	pg_queries "dbdaddy/db/pg/queries"
)

func GetExistingDbs() ([]string, error) {
	rows, err := db.DB.Query(pg_queries.QGetExistingDbs())
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
	_, err := db.DB.Query(pg_queries.QDisconnectAllUsersFromDb(dbname))
	return err
}

func DbExists(dbname string) bool {
	row := db.DB.QueryRow(pg_queries.QCheckDbExists(dbname))

	exists := false
	_ = row.Scan(&exists)

	return exists
}

func NewDbFromOriginal(originalDbName string, newDbName string) error {
	err := DisconnectAllUsers(originalDbName)
	if err != nil {
		return err
	}

	_, err = db.DB.Query(pg_queries.QCreateNewDbFromOldTemplate(newDbName, originalDbName))
	if err != nil {
		return err
	}

	return nil
}

func CreateDb(dbname string) error {
	_, err := db.DB.Query(pg_queries.QCreateNewDb(dbname))
	return err
}

func DeleteDb(dbname string) error {
	if err := DisconnectAllUsers(dbname); err != nil {
		return err
	}

	_, err := db.DB.Query(pg_queries.QDeleteDb(dbname))
	if err != nil {
		return err
	}

	return nil
}
