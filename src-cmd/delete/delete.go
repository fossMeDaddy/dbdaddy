package deleteCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/middlewares"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "delete <branchname>",
	Short: "Deletes a database branch, this action can not be undone",
	Args:  cobra.ExactArgs(1),
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	branchname := args[0]

	if !db_int.DbExists(branchname) {
		cmd.PrintErrf("Database branch '%s' does not exist\n", branchname)
		return
	}

	err := db_int.DeleteDb(branchname)
	if err != nil {
		cmd.PrintErrln("Error deleting database\n" + err.Error())
		return
	}

	cmd.Printf("Successfully deleted db '%s'\n", branchname)

	if viper.GetString(constants.DbConfigCurrentBranchKey) == branchname {
		existingDbs, err := db_int.GetExistingDbs()
		if err != nil {
			panic("Something went wrong!\n" + err.Error())
		}

		newBranchName := existingDbs[0]
		viper.Set(constants.DbConfigCurrentBranchKey, newBranchName)
		viper.WriteConfig()

		cmd.Printf("Switched to database branch '%s'\n", newBranchName)
	}
}

func Init() *cobra.Command {
	// add some flags
	// ...

	return cmd
}
