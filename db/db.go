package db

import (
	"database/sql"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/globals"

	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

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

	db, err := ConnectDb(v, constants.SelfDbName)
	if errors.Is(err, errs.ErrUnsupportedDriver) {
		return db, err
	} else if err != nil {
		userDb, userErr := ConnectDb(v, v.GetString(constants.DbConfigDbNameKey))
		if userErr != nil {
			return nil, fmt.Errorf("error connecting to your database!\n" + err.Error())
		}
		defer userDb.Close()

		if _, err := userDb.Query(fmt.Sprintf("CREATE DATABASE %s", constants.SelfDbName)); err != nil {
			return nil, fmt.Errorf("could not create database, please check your connection or permissions!\n" + err.Error())
		}

		selfDb, selfErr := ConnectDb(v, constants.SelfDbName)
		if selfErr != nil {
			return nil, fmt.Errorf("unexpected error occured!\n" + selfErr.Error())
		}

		fnLocalDb = selfDb
	} else {
		fnLocalDb = db
	}

	globals.ConnDbName = constants.SelfDbName
	globals.DB = fnLocalDb

	return globals.DB, nil
}

// deprecated in favor of makeing things db driver agnostic
func ConnectSelfDb_DEPRECATED(v *viper.Viper) (*sql.DB, error) {
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
	var (
		db  *sql.DB
		err error
	)

	switch v.GetString(constants.DbConfigDriverKey) {
	case constants.DbDriverPostgres:
		db, err = openConn("pgx", GetPgConnUriFromViper(v, dbname))
	case constants.DbDriverMySQL:
		db, err = openConn("mysql", GetMysqlConnUriFromViper(v, dbname))
	default:
		return db, errs.ErrUnsupportedDriver
	}

	if err != nil {
		return db, err
	}

	globals.ConnDbName = dbname
	globals.DB = db

	return globals.DB, nil
}
