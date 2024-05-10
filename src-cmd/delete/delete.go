package deleteCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/pg"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "delete <branchname>",
	Short: "Deletes a database branch, this action can not be undone",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	branchname := args[0]

	if branchname == constants.SelfDbName {
		cmd.PrintErrf("Error deleting database branch '%s', bro seriously... why'd you do that?\n", branchname)
		return
	}

	if !pg.DbExists(branchname) {
		cmd.PrintErrf("Database branch '%s' does not exist\n", branchname)
		return
	}

	err := pg.DeleteDb(branchname)
	if err != nil {
		cmd.PrintErrln("Error deleting database\n" + err.Error())
		return
	}

	cmd.Printf("Successfully deleted db '%s'\n", branchname)

	if viper.GetString(constants.DbConfigCurrentBranchKey) == branchname {
		existingDbs, err := pg.GetExistingDbs()
		if err != nil {
			panic("Something went wrong!\n" + err.Error())
		}

		fmt.Println(existingDbs)
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
