package middlewares

import (
	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckConnection(fn types.CobraCmdFn) types.CobraCmdFn {
	return func(cmd *cobra.Command, args []string) {
		_, err := db.ConnectSelfDb(viper.GetViper())
		if err != nil {
			cmd.PrintErrln(err.Error())
			cmd.PrintErrln("\nTry reviewing your database credentials in the config file, if not present, try creating one. (hint: 'config -h')")
			return
		}

		fn(cmd, args)
	}
}
