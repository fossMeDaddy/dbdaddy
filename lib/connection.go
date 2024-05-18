package lib

import (
	"dbdaddy/db"
	"fmt"
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
