package remoteAddCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var (
	forceFlag bool
)

var cmd = &cobra.Command{
	Use:     "add",
	Short:   "add an origin by a providing a name and a db connection uri for the origin",
	Example: "add origin postgresql://user:pwd@localhost:5432/dbname",
	Run:     run,
	Args:    cobra.ExactArgs(2),
}

func run(cmd *cobra.Command, args []string) {
	originName := args[0]
	connUri := args[1]

	connConfig, uriErr := db_int.GetConnConfigFromUri(connUri)
	if uriErr != nil {
		cmd.PrintErrln(uriErr)
		return
	}

	originConfigKey := libUtils.GetDbConfigOriginKey(originName)

	origin := viper.GetStringMapString(originConfigKey)
	if len(maps.Keys(origin)) > 0 && !forceFlag {
		cmd.Println(fmt.Sprintf("remote origin '%s' already exists in the config, use force flag to override", originName))
		return
	}

	viper.Set(originConfigKey, connConfig)
	err := viper.WriteConfig()
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("Remote origin '%s' successfully set.", originName))
}

func Init() *cobra.Command {
	// flags here
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "will override existing configuration for provided origin")

	return cmd
}
