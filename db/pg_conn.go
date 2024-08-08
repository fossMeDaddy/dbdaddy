package db

import (
	constants "dbdaddy/const"
	"fmt"

	"github.com/spf13/viper"
)

func GetPgConnUriFromViper(v *viper.Viper, dbname string) string {
	// EXAMPLE URI:
	// postgresql://sally:sallyspassword@dbserver.example:5555/userdata?connect_timeout=10&sslmode=require&target_session_attrs=primary

	host := v.GetString(constants.DbConfigHostKey)
	port := v.GetUint16(constants.DbConfigPortKey)
	user := v.GetString(constants.DbConfigUserKey)
	password := v.GetString(constants.DbConfigPassKey)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, fmt.Sprint(port), dbname)
}
