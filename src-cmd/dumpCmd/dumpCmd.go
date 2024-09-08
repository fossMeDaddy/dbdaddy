package dumpCmd

import (
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/fossmedaddy/dbdaddy/src-cmd/dumpCmd/dumpCleanCmd"
	"github.com/fossmedaddy/dbdaddy/src-cmd/dumpCmd/dumpListCmd"

	"github.com/spf13/cobra"
)

var (
	useGlobalConfigFile bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "dumps",
	Short: "perform operations on db backups aka \"dumps\" i.e. list or clear them",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Init() *cobra.Command {
	cmd.AddCommand(dumpListCmd.Init())
	cmd.AddCommand(dumpCleanCmd.Init())

	cmd.Flags().BoolVarP(&useGlobalConfigFile, "global", "g", false, "explicitly use the global config file to list the dumps")

	return cmd
}
