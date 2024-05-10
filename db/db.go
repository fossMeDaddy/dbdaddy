package db

import (
	"database/sql"
	constants "dbdaddy/const"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

var DB *sql.DB

var SupportedDrivers = []string{"postgres"}

func getPgConnUriFromViper(dbname string) string {
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
		db, err := OpenConnection("pgx", getPgConnUriFromViper(constants.SelfDbName))
		if err != nil {
			userDb, userErr := OpenConnection("pgx", getPgConnUriFromViper(viper.GetString(constants.DbConfigDbNameKey)))
			if userErr != nil {
				return nil, fmt.Errorf("error connecting to your database!\n" + err.Error())
			}

			if _, err := userDb.Query(fmt.Sprintf("CREATE DATABASE %s", constants.SelfDbName)); err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					return nil, fmt.Errorf("could not create database, please check your connection!\n" + err.Error())
				}
			}

			db, err := OpenConnection("pgx", getPgConnUriFromViper(constants.SelfDbName))
			if err != nil {
				return nil, fmt.Errorf("unexpected error occured!\n" + err.Error())
			}

			DB = db

			userDb.Close()
		} else {
			DB = db
		}

		return DB, nil
	}

	panic(fmt.Sprintf("'%s' driver is not supported, as of now, the supported drivers are: %v", viper.Get(constants.DbConfigDriverKey), SupportedDrivers))
}
