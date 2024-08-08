package msql

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/db/msql/msqlq"
)

func GetExistingDbs() ([]string, error) {
	rows, err := db.DB.Query(msqlq.QGetExistingDbs())
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
