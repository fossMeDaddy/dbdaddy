package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/globals"

	_ "github.com/jackc/pgx/v5/stdlib"
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

func openConn(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping unsucessful: %s", err.Error())
	}

	return db, nil
}

func ConnectSelfDb(v *viper.Viper) (*sql.DB, error) {
	var fnLocalDb *sql.DB
	if viper.Get(constants.DbConfigDriverKey) == constants.DbDriverPostgres {
		db, err := openConn("pgx", GetPgConnUriFromViper(v, constants.SelfDbName))
		if err != nil {
			userDb, userErr := openConn("pgx", GetPgConnUriFromViper(v, v.GetString(constants.DbConfigDbNameKey)))
			if userErr != nil {
				return nil, fmt.Errorf("error connecting to your database!\n" + err.Error())
			}

			if _, err := userDb.Exec(fmt.Sprintf("CREATE DATABASE %s", constants.SelfDbName)); err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					return nil, fmt.Errorf("could not create database, please check your connection!\n" + err.Error())
				}
			}

			db, err := openConn("pgx", GetPgConnUriFromViper(v, constants.SelfDbName))
			if err != nil {
				return nil, fmt.Errorf("unexpected error occured!\n" + err.Error())
			}

			userDb.Close()

			fnLocalDb = db
		} else {
			fnLocalDb = db
		}

		globals.ConnDbName = constants.SelfDbName
		globals.DB = fnLocalDb

		return globals.DB, nil
	}

	fmt.Println(viper.Get(constants.DbConfigDriverKey), constants.SupportedDrivers)
	panic(fmt.Sprintf("'%s' driver is not supported, as of now, the supported drivers are: %v", viper.Get(constants.DbConfigDriverKey), constants.SupportedDrivers))
}

func ConnectDb(v *viper.Viper, dbname string) (*sql.DB, error) {
	if v.GetString(constants.DbConfigDriverKey) == constants.DbDriverPostgres {
		db, err := openConn("pgx", GetPgConnUriFromViper(v, dbname))
		if err != nil {
			return nil, fmt.Errorf("unexpected error occured!\n" + err.Error())
		}

		globals.ConnDbName = dbname
		globals.DB = db

		return globals.DB, nil
	}

	fmt.Println(v.Get(constants.DbConfigDriverKey), constants.SupportedDrivers)
	panic(fmt.Sprintf("'%s' driver is not supported, as of now, the supported drivers are: %v", v.Get(constants.DbConfigDriverKey), constants.SupportedDrivers))
}
