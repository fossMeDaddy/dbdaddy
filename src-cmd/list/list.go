package listCmd

import (
	db_lib "dbdaddy/db/lib"
	"dbdaddy/middlewares"
	"fmt"

	"github.com/spf13/cobra"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "list",
	Short: "Lists available databases on the server",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	dbs, err := db_lib.GetExistingDbs()
	if err != nil {
		cmd.PrintErrln("Unexpected error occured!\n" + err.Error())
	}

	cmd.Println("Available databases:")
	for i, db := range dbs {
		cmd.Println(fmt.Sprintf("%d. %s", i+1, db))
	}
}

func Init() *cobra.Command {
	// add flags

	return cmd
}
