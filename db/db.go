package db

import (
	"database/sql"
	constants "dbdaddy/const"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

var DB *sql.DB

func GetPgConnUriFromViper(dbname string) string {
	// EXAMPLE URI:
	// postgresql://sally:sallyspassword@dbserver.example:5555/userdata?connect_timeout=10&sslmode=require&target_session_attrs=primary

	host := viper.GetString(constants.DbConfigHostKey)
	port := viper.GetUint16(constants.DbConfigPortKey)
	user := viper.GetString(constants.DbConfigUserKey)
	password := viper.GetString(constants.DbConfigPassKey)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, fmt.Sprint(port), dbname)
}

func OpenConnection(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping unsucessful: %s", err.Error())
	}

	return db, nil
}

func ConnectDB() (*sql.DB, error) {
	if viper.Get(constants.DbConfigDriverKey) == constants.DbDriverPostgres {
		currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)
		db, err := OpenConnection("pgx", GetPgConnUriFromViper(currBranch))
		if err != nil {
			return nil, fmt.Errorf("error connecting to your database!\n" + err.Error())
		} else {
			DB = db
		}

		return DB, nil
	}

	fmt.Println(viper.Get(constants.DbConfigDriverKey), constants.SupportedDrivers)
	panic(fmt.Sprintf("'%s' driver is not supported, as of now, the supported drivers are: %v", viper.Get(constants.DbConfigDriverKey), constants.SupportedDrivers))
}
