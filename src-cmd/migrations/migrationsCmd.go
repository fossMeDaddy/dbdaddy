package migrationsCmd

import (
	"github.com/fossmedaddy/dbdaddy/src-cmd/migrations/migrationsDown"
	migrationsGenCmd "github.com/fossmedaddy/dbdaddy/src-cmd/migrations/migrationsGen"
	migrationsStatusCmd "github.com/fossmedaddy/dbdaddy/src-cmd/migrations/migrationsStatus"
	"github.com/fossmedaddy/dbdaddy/src-cmd/migrations/migrationsUp"

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

	// commands
	cmd.AddCommand(migrationsGenCmd.Init())
	cmd.AddCommand(migrationsStatusCmd.Init())
	cmd.AddCommand(migrationsUp.Init())
	cmd.AddCommand(migrationsDown.Init())

	return cmd
}
