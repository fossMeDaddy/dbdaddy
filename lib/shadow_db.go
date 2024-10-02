package lib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
)

// NOTE: globals.DB should be connected to the shadow database before running this function
func CreateShadowDB() (string, error) {
	shadowDbPrefix := fmt.Sprintf("%s_", constants.ShadowDbPrefix)
	seq := -1

	dbs, err := db_int.GetExistingDbs(true)
	if err != nil {
		return "", err
	}

	for _, db := range dbs {
		if strings.HasPrefix(shadowDbPrefix, db) {
			dbSeq, err := strconv.Atoi(strings.ReplaceAll(db, shadowDbPrefix, ""))
			if err != nil {
				return "", fmt.Errorf("database name '%s' has invalid non-parsable sequence, please delete this db to resolve the error (author skill issues detected).")
			}

			seq = max(dbSeq, seq)
		}
	}

	seq++

	shadowDbName := shadowDbPrefix + fmt.Sprint(seq)
	if err := db_int.CreateDb(shadowDbName); err != nil {
		return "", err
	}

	return shadowDbName, err
}
