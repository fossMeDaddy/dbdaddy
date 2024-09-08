package deleteCmd

import (
	"dbdaddy/constants"
	"dbdaddy/db/db_int"
	"dbdaddy/middlewares"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	silentFlag bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "delete <branchname>",
	Short: "Deletes a database branch, this action can not be undone",
	Args:  cobra.MaximumNArgs(1),
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	var branchname string
	if len(args) > 0 {
		branchname = args[0]

		if !db_int.DbExists(branchname) {
			cmd.PrintErrf("Database branch '%s' does not exist\n", branchname)
			return
		}
	} else {
		dbs, err := db_int.GetExistingDbs()
		if err != nil {
			cmd.Println("unexpected error occured", err)
			return
		}

		dbPrompt := promptui.Select{
			Label: "Choose database to delete",
			Items: dbs,
			Searcher: func(input string, index int) bool {
				return strings.Contains(dbs[index], strings.ToLower(strings.Trim(input, " ")))
			},
			StartInSearchMode: true,
		}

		_, result, pErr := dbPrompt.Run()
		if pErr != nil {
			cmd.PrintErrln("Error occured while choosing deletion database", pErr)
			return
		}

		branchname = result
	}

	if !silentFlag {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Deleting database: '%s' continue", branchname),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			cmd.Println("Cancelling deletion.")
			return
		}
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
	cmd.Flags().BoolVar(&silentFlag, "silent", false, "silences confirmation for deleting branch")

	return cmd
}
