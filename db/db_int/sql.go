package db_int

import (
	constants "dbdaddy/const"

	"github.com/spf13/viper"
)

func GetDisableConstSQL() string {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL DEFERRED;"
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 0;"
	default:
		panic("")
	}
}

func GetEnableConstSQL() string {
	driver := viper.GetString(constants.DbConfigDriverKey)
	switch driver {
	case constants.DbDriverPostgres:
		return "SET CONSTRAINTS ALL IMMEDIATE;"
	case constants.DbDriverMySQL:
		return "SET FOREIGN_KEY_CHECKS = 1;"
	default:
		panic("")
	}
}
