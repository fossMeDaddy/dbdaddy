package db

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/spf13/viper"
)

func GetMysqlConnUriFromViper(v *viper.Viper, dbname string) string {
	return fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?multiStatements=true`,
		v.GetString(constants.DbConfigUserKey),
		v.GetString(constants.DbConfigPassKey),
		v.GetString(constants.DbConfigHostKey),
		v.GetString(constants.DbConfigPortKey),
		dbname,
	)
}
