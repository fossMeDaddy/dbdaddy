package lib

import (
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"

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
