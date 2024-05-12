package checkoutCmd

import (
	constants "dbdaddy/const"
	db_lib "dbdaddy/db/lib"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags
var (
	shouldCreateNewBranch bool
)

var cmd = &cobra.Command{
	Use:   "checkout <branchname>",
	Short: "checkout into a new/existing branch in the database",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func Init() *cobra.Command {
	// '-b' creates a new branch with a unique name
	cmd.Flags().BoolVarP(&shouldCreateNewBranch, "new", "n", false, "create a new branch with given branch name")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	branchname := args[0]

	if shouldCreateNewBranch {
		err := db_lib.NewDbFromOriginal(viper.GetString(constants.DbConfigCurrentBranchKey), branchname)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				cmd.PrintErrf("Could not create a new database branch with name '%s' because it already exists.\n", branchname)
				return
			}

			panic(err)
		}
	} else {
		if db_lib.DbExists(branchname) {
			viper.Set(constants.DbConfigCurrentBranchKey, branchname)
			viper.WriteConfig()
		} else {
			cmd.PrintErrf("Database branch '%s' does not exists, run 'checkout <branchname> -n' to create a new branch\n", branchname)
			return
		}
	}

	viper.Set(constants.DbConfigCurrentBranchKey, branchname)
	viper.WriteConfig()

	cmd.Println(fmt.Sprintf("Switched to branch: %s", branchname))
}
