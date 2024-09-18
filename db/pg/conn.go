package pg

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/jackc/pgx/v5"
)

func GetConnConfigFromUri(uri string) (types.ConnConfig, error) {
	connConfig := types.ConnConfig{}

	dbConfig, uriErr := pgx.ParseConfig(uri)
	if uriErr != nil {
		return connConfig, uriErr
	}

	connConfig.User = dbConfig.User
	connConfig.Password = dbConfig.Password
	connConfig.Host = dbConfig.Host
	connConfig.Port = fmt.Sprint(dbConfig.Port)
	connConfig.Database = dbConfig.Database
	connConfig.Params = dbConfig.RuntimeParams
	connConfig.Driver = constants.DbDriverPostgres

	return connConfig, nil
}
