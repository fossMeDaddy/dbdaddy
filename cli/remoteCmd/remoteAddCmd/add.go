package remoteAddCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

const cmdManual = `
Name of the origin is same as the database name the origin is being set up of for.
`

var (
	forceFlag bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:     "add",
	Short:   "add a remote origin by a providing a name and a db connection uri for the current branch",
	Example: "add postgresql://user:pwd@localhost:5432/dbname",
	Run:     cmdRunFn,
	Args:    cobra.ExactArgs(1),
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	connUri := args[0]

	originConfigKey := libUtils.GetDbConfigOriginKey(currBranch)

	origin := viper.GetStringMap(originConfigKey)
	if len(maps.Keys(origin)) > 0 && !forceFlag {
		cmd.Println(fmt.Sprintf("remote origin for '%s' already exists in the config, use force flag to override", currBranch))
		return
	}

	connConfig, uriErr := libUtils.GetConnConfigFromUri(connUri)
	if uriErr != nil {
		cmd.PrintErrln(uriErr)
		return
	}

	if err := lib.TmpSwitchConn(connConfig, func() error {
		return nil
	}); err != nil {
		cmd.PrintErrln(err)
		return
	}

	viper.Set(originConfigKey, connConfig)
	err := viper.WriteConfig()
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("remote origin for '%s' successfully set.", currBranch))
}

func Init() *cobra.Command {
	// flags here
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "will override existing configuration for provided origin")

	return cmd
}
