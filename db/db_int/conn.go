package db_int

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/pg"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/spf13/viper"
)

func GetConnConfigFromUri(uri string) (types.ConnConfig, error) {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.GetConnConfigFromUri(uri)
	default:
		return types.ConnConfig{}, errs.ErrUnsupportedDriver
	}
}
