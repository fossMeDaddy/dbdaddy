package remoteCmd

import (
	"github.com/fossmedaddy/dbdaddy/cli/remoteCmd/remoteAddCmd"
	"github.com/fossmedaddy/dbdaddy/cli/remoteCmd/remoteLsCmd"
	"github.com/fossmedaddy/dbdaddy/cli/remoteCmd/remoteRemoveCmd"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "remote",
	Short: "add or remove remote origins",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Init() *cobra.Command {
	// flags here

	// commands
	cmd.AddCommand(remoteAddCmd.Init())
	cmd.AddCommand(remoteLsCmd.Init())
	cmd.AddCommand(remoteRemoveCmd.Init())

	return cmd
}
