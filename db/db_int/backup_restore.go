package db_int

import (
	constants "dbdaddy/const"
	"dbdaddy/db/msql"
	"dbdaddy/db/pg"
	"dbdaddy/errs"

	"github.com/spf13/viper"
)

func DumpDb(outputFilePath string, v *viper.Viper) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.DumpDb(outputFilePath, v)
	case constants.DbDriverMySQL:
		return msql.DumpDb(outputFilePath, v)
	default:
		return errs.ErrUnsupportedDriver
	}
}
func DumpDbOnlySchema(outputFilePath string, v *viper.Viper) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.DumpDbOnlySchema(outputFilePath, v)
	} else {
		panic("")
	}
}

func RestoreDb(dbname string, v *viper.Viper, dumpFilePath string, override bool) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return pg.RestoreDb(dbname, v, dumpFilePath, override)
	case constants.DbDriverMySQL:
		return msql.RestoreDb(dbname, v, dumpFilePath, override)
	default:
		return errs.ErrUnsupportedDriver
	}
}
func RestoreDbOnlySchema(dbname string, v *viper.Viper, dumpFilePath string) error {
	driver := v.GetString(constants.DbConfigDriverKey)
	if driver == constants.DbDriverPostgres {
		return pg.RestoreDbOnlySchema(dbname, v, dumpFilePath)
	} else {
		panic("")
	}
}
