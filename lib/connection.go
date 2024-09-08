package lib

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/db"

	"github.com/spf13/viper"
)

func PingDB() error {
	if db.DB == nil {
		return fmt.Errorf("error connecting to database, db object is <nil>")
	}

	err := db.DB.Ping()
	if err != nil {
		return fmt.Errorf("error occured while pinging database: %s", err.Error())
	}

	return nil
}

// error can be returned if callback function errors or there is a connection error to the DB
func SwitchDB(v *viper.Viper, dbname string, fn func() error) error {
	defer db.ConnectDb(viper.GetViper(), db.ConnDbName)

	_, err := db.ConnectDb(v, dbname)
	if err != nil {
		return err
	}

	fnErr := fn()
	if fnErr != nil {
		return fnErr
	}

	return nil
}
