package libUtils

import (
	"fmt"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func GetConnConfigFromViper(v *viper.Viper) (types.ConnConfig, error) {
	connConfig := types.ConnConfig{}
	parseErr := v.UnmarshalKey(constants.DbConfigConnKey, &connConfig)
	if parseErr != nil {
		return connConfig, parseErr
	}

	currBranch := v.GetString(constants.DbConfigCurrentBranchKey)
	connConfig.Database = currBranch

	return connConfig, nil
}

func GetConnConfigFromUri(uri string) (types.ConnConfig, error) {
	connConfig := types.ConnConfig{}

	if strings.HasPrefix(uri, "postgresql://") {
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
	} else {
		return connConfig, errs.ErrUnsupportedDriver
	}

	return connConfig, nil
}

func GetShadowConnConfig(v *viper.Viper) (types.ConnConfig, error) {
	connConfig := types.ConnConfig{}

	if v.IsSet(constants.DbConfigShadowConnKey) {
		if err := v.UnmarshalKey(constants.DbConfigShadowConnKey, &connConfig); err != nil {
			return connConfig, err
		}
	} else {
		if cc, err := GetConnConfigFromViper(v); err != nil {
			return cc, err
		} else {
			connConfig = cc
		}
	}

	return connConfig, nil

}
