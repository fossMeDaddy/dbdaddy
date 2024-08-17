package migrationsCmd

import (
	migrationsGenCmd "dbdaddy/migrationsCmd/migrationsGen"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "migrations",
	Short: "manage migrations, see 'migrations --help' for more info.",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Init() *cobra.Command {
	// flags

	cmd.AddCommand(migrationsGenCmd.Init())

	return cmd
}
