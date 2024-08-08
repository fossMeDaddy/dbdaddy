package db

import (
	constants "dbdaddy/const"
	"fmt"

	"github.com/spf13/viper"
)

func GetMysqlConnUriFromViper(v *viper.Viper, dbname string) string {
	return fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s`,
		v.GetString(constants.DbConfigUserKey),
		v.GetString(constants.DbConfigPassKey),
		v.GetString(constants.DbConfigHostKey),
		v.GetString(constants.DbConfigPortKey),
		dbname,
	)
}
