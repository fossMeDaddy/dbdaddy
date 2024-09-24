package remoteRemoveCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:     "remove",
	Short:   "remove origin(s) by name",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Run:     run,
}

func run(cmd *cobra.Command, args []string) {
	origins := types.DbConfigOrigins{}
	parseErr := viper.UnmarshalKey(constants.DbConfigOriginsKey, &origins)

	if parseErr != nil {
		cmd.PrintErrln(parseErr)
		return
	}

	for _, argOrigin := range args {
		delete(origins, argOrigin)
	}

	viper.Set(constants.DbConfigOriginsKey, origins)
	err := viper.WriteConfig()
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("Removed origins %v successfully.", args))
}

func Init() *cobra.Command {
	// flags here

	return cmd
}
