package pg

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/pg/pgq"
	"github.com/fossmedaddy/dbdaddy/db/sharedq"
	"github.com/fossmedaddy/dbdaddy/globals"
)

func GetExistingDbs() ([]string, error) {
	rows, err := globals.DB.Query(pgq.QGetExistingDbs())
	if err != nil {
		return nil, err
	}

	existingDbs := []string{}
	defer rows.Close()
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
	_, err := globals.DB.Exec(pgq.QDisconnectAllUsersFromDb(dbname))
	return err
}

func DbExists(dbname string) bool {
	row := globals.DB.QueryRow(pgq.QCheckDbExists(dbname))

	exists := false
	_ = row.Scan(&exists)

	return exists
}

// DO NOT USE, disconnects all users while creating new db
func NewDbFromOriginal_DEPRECATED(originalDbName string, newDbName string) error {
	err := DisconnectAllUsers(originalDbName)
	if err != nil {
		return err
	}

	_, err = globals.DB.Exec(pgq.QCreateNewDbFromOldTemplate(newDbName, originalDbName))
	if err != nil {
		return err
	}

	return nil
}

func CreateDb(dbname string) error {
	_, err := globals.DB.Exec(sharedq.QCreateNewDb(dbname))
	return err
}

func DeleteDb(dbname string) error {
	if err := DisconnectAllUsers(dbname); err != nil {
		return err
	}

	_, err := globals.DB.Exec(sharedq.QDeleteDb(dbname))
	if err != nil {
		return err
	}

	return nil
}
