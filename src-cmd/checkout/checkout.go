package checkoutCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/middlewares"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags
var (
	shouldCreateNewBranch bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "checkout <branchname>",
	Short: "checkout into a new/existing branch in the database",
	Args:  cobra.ExactArgs(1),
	Run:   cmdRunFn,
}

func Init() *cobra.Command {
	// '-b' creates a new branch with a unique name
	cmd.Flags().BoolVarP(&shouldCreateNewBranch, "new", "n", false, "create a new branch with given branch name")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	branchname := args[0]

	if shouldCreateNewBranch {
		err := db_int.NewDbFromOriginal(viper.GetString(constants.DbConfigCurrentBranchKey), branchname)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				cmd.PrintErrf("Could not create a new database branch with name '%s' because it already exists.\n", branchname)
				return
			}

			panic(err)
		}
	} else {
		if db_int.DbExists(branchname) {
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
