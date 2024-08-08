package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/db/pg/pgq"
)

func GetExistingDbs() ([]string, error) {
	rows, err := db.DB.Query(pgq.QGetExistingDbs())
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
	_, err := db.DB.Query(pgq.QDisconnectAllUsersFromDb(dbname))
	return err
}

func DbExists(dbname string) bool {
	row := db.DB.QueryRow(pgq.QCheckDbExists(dbname))

	exists := false
	_ = row.Scan(&exists)

	return exists
}

func NewDbFromOriginal(originalDbName string, newDbName string) error {
	err := DisconnectAllUsers(originalDbName)
	if err != nil {
		return err
	}

	_, err = db.DB.Query(pgq.QCreateNewDbFromOldTemplate(newDbName, originalDbName))
	if err != nil {
		return err
	}

	return nil
}

func CreateDb(dbname string) error {
	_, err := db.DB.Query(pgq.QCreateNewDb(dbname))
	return err
}

func DeleteDb(dbname string) error {
	if err := DisconnectAllUsers(dbname); err != nil {
		return err
	}

	_, err := db.DB.Query(pgq.QDeleteDb(dbname))
	if err != nil {
		return err
	}

	return nil
}
