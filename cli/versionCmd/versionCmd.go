package versionCmd

import (
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "version",
	Short: "self-explanatory",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Println(globals.Version)
}

func Init() *cobra.Command {
	// add flags here

	return cmd
}
