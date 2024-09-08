package pg

import (
	"dbdaddy/constants"
	"fmt"

	"github.com/spf13/viper"
)

func GetPgConnUriFromViper(dbname string) string {
	// EXAMPLE URI:
	// postgresql://sally:sallyspassword@dbserver.example:5555/userdata?connect_timeout=10&sslmode=require&target_session_attrs=primary

	host := viper.GetString(constants.DbConfigHostKey)
	port := viper.GetUint16(constants.DbConfigPortKey)
	user := viper.GetString(constants.DbConfigUserKey)
	password := viper.GetString(constants.DbConfigPassKey)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, fmt.Sprint(port), dbname)
}
