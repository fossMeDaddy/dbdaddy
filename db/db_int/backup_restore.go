package db_int

import (
	constants "dbdaddy/const"
	"dbdaddy/db/pg"

	"github.com/spf13/viper"
)

func TakeADump(outputFilePath string, v *viper.Viper) error {
	driver := viper.GetString(constants.DbConfigDriverKey)
	if driver == viper.GetString(constants.DbConfigDriverKey) {
		return pg.TakeADump(outputFilePath, v)
	} else {
		panic("")
	}
}
