package lib

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/types"
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

// takes in the currently active connection config to switch back to, dbname and
// function to run when connected to the "dbname" database.
// error can be returned if callback function errors or there is a connection error to the DB
func TmpSwitchDB(dbname string, fn func() error) error {
	if globals.DB == nil {
		return fmt.Errorf("database not connected, connect before switching between databases, author skill issues...")
	}

	defer db.ConnectDb(globals.CurrentConnConfig)

	connConfig := globals.CurrentConnConfig
	connConfig.Database = dbname
	_, err := db.ConnectDb(connConfig)
	if err != nil {
		return err
	}

	return fn()
}

func TmpSwitchConn(connConfig types.ConnConfig, fn func() error) error {
	if globals.DB == nil {
		return fmt.Errorf("database not connected, should not be happening here, author akill issues...")
	}

	defer db.ConnectDb(globals.CurrentConnConfig)

	_, err := db.ConnectDb(connConfig)
	if err != nil {
		return err
	}

	return fn()
}
