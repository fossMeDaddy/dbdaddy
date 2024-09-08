package lib

import (
	"dbdaddy/constants"
	"dbdaddy/lib/libUtils"
	"path"

	"github.com/spf13/viper"
)

func GetDriverDumpDir(configPath string) string {
	v := viper.New()
	ReadConfig(v, configPath)

	configDirPath, _ := path.Split(configPath)

	dumpDirName := constants.DriverDumpDirNames[v.GetString(constants.DbConfigDriverKey)]

	dumpDirPath := path.Join(configDirPath, dumpDirName)
	libUtils.EnsureDirExists(dumpDirPath)

	return dumpDirPath
}
