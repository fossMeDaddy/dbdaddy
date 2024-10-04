package msql

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/msql/msqlq"
	"github.com/fossmedaddy/dbdaddy/db/sharedq"
	"github.com/fossmedaddy/dbdaddy/globals"
)

func GetExistingDbs(showHidden bool) ([]string, error) {
	rows, err := globals.DB.Query(msqlq.QGetExistingDbs())
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

func DbExists(dbname string) bool {
	row := globals.DB.QueryRow(msqlq.QCheckDbExists(dbname))
	if row.Err() != nil {
		return false
	}

	dbCount := 0
	err := row.Scan(&dbCount)
	if err != nil {
		return false
	}

	return dbCount == 1
}

func CreateDb(dbname string) error {
	_, err := globals.DB.Query(sharedq.QCreateNewDb(dbname))
	return err
}

func DeleteDb(dbname string) error {
	_, err := globals.DB.Query(sharedq.QDeleteDb(dbname))
	if err != nil {
		return err
	}

	return nil
}
