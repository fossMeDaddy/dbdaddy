package lib

import (
	constants "dbdaddy/const"
	"path"

	"github.com/spf13/viper"
)

func GetDriverDumpDir(configPath string) string {
	v := viper.New()
	ReadConfig(v, configPath)

	configDirPath, _ := path.Split(configPath)

	dumpDirName := constants.DriverDumpDirNames[v.GetString(constants.DbConfigDriverKey)]

	dumpDirPath := path.Join(configDirPath, dumpDirName)
	DirExistsCreate(dumpDirPath)

	return dumpDirPath
}
