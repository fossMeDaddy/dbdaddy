package middlewares

import (
	"dbdaddy/db"
	"dbdaddy/lib"
	"dbdaddy/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckConnection(fn types.CobraCmdFn) types.CobraCmdFn {
	return func(cmd *cobra.Command, args []string) {
		configFilePath, _ := lib.FindConfigFilePath()
		lib.InitConfigFile(viper.GetViper(), configFilePath, false)

		_, err := db.ConnectDB()
		if err != nil {
			cmd.PrintErrln(err.Error())
			cmd.PrintErrln("\nTry reviewing your database credentials in the config file, if not present, try creating one. (hint: 'config -h')")
			return
		}

		fn(cmd, args)
	}
}
