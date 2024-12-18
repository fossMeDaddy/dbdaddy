package listCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/middlewares"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showHiddenFlag bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "Lists available database branches on the db server",
	Run:     cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	dbs, err := db_int.GetExistingDbs(showHiddenFlag)
	if err != nil {
		cmd.PrintErrln("Unexpected error occured!\n" + err.Error())
		return
	}

	cmd.Println("Available database branches:")
	for i, db := range dbs {
		dbStr := fmt.Sprintf("%d - %s", i+1, db)
		if db == currBranch {
			dbStr += " (current branch)"
		}

		cmd.Println(dbStr)
	}
}

func Init() *cobra.Command {
	// add flags
	cmd.Flags().BoolVarP(&showHiddenFlag, "show-hidden", "s", false, "Show hidden databases")

	return cmd
}
