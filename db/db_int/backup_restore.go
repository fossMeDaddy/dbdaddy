package db_int

import (
	constants "dbdaddy/const"
	"dbdaddy/db/pg"

	"github.com/spf13/viper"
)

func DumpDb(outputFilePath string, v *viper.Viper) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DumpDb(outputFilePath, v)
	} else {
		panic("")
	}
}

func RestoreDb(dbname string, v *viper.Viper, dumpFilePath string, override bool) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.RestoreDb(dbname, v, dumpFilePath, override)
	} else {
		panic("")
	}
}
