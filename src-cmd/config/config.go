package configCmd

import (
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "config",
	Short: "Used to manage config files, check subcommands by running 'config -h'",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Init() *cobra.Command {
	cmd.AddCommand(InitCreateCmd())
	cmd.AddCommand(InitOpenCmd())

	return cmd
}
