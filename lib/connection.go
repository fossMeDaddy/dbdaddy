package lib

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

func PingDB() error {
	if globals.DB == nil {
		return fmt.Errorf("error connecting to database, db object is <nil>")
	}

	err := globals.DB.Ping()
	if err != nil {
		return fmt.Errorf("error occured while pinging database: %s", err.Error())
	}

	return nil
}

// error can be returned if callback function errors or there is a connection error to the DB
func SwitchDB(v *viper.Viper, dbname string, fn func() error) error {
	prevConnConfig := types.ConnConfig{}
	if err := v.UnmarshalKey(constants.DbConfigConnSubkey, &prevConnConfig); err != nil {
		return err
	}
	prevConnConfig.Database = globals.ConnDbName

	defer db.ConnectDb(prevConnConfig)

	connConfig := prevConnConfig
	connConfig.Database = dbname
	_, err := db.ConnectDb(connConfig)
	if err != nil {
		return err
	}

	fnErr := fn()
	if fnErr != nil {
		return fnErr
	}

	return nil
}
