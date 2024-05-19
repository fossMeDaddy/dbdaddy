package restoreCmd

import (
	"dbdaddy/middlewares"

	"github.com/spf13/cobra"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "restore <dbname>",
	Short: "restores the given database with a database dump",
	Long:  "restoration defaults to restoring the db with the latest dump taken. if no dumps are present with me & i'll require you to provide a path to the backed up dump",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
}

func Init() *cobra.Command {
	return cmd
}
