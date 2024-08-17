package checkoutCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/errs"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags
var (
	shouldCreateNewBranch bool
	shouldKeepItClean     bool
	shouldCopyOnlySchema  bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:     "checkout <branchname>",
	Aliases: []string{"checkoutmedaddy"},
	Short:   "checkout into a new/existing branch in the database",
	Args:    cobra.ExactArgs(1),
	Run:     cmdRunFn,
}

func Init() *cobra.Command {
	cmd.Flags().BoolVarP(&shouldCreateNewBranch, "new", "n", false, "create a new branch with given branch name, the contents of the current branch will be copied over schema and data both")
	cmd.Flags().BoolVarP(&shouldKeepItClean, "clean", "c", false, "the new branch created by '-n' will be independent of the current branch i.e. nothing will be copied over")
	cmd.Flags().BoolVar(&shouldCopyOnlySchema, "only-schema", false, "create a new branch with given branch name, only the schema of the current branch will be copied over")
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	branchname := args[0]
	// flags validation
	if shouldKeepItClean && !shouldCreateNewBranch {
		cmd.PrintErrln("'--clean, -c' option can only be provided when creating new branches i.e. with '-n, --new' option")
		return
	}
	if shouldCopyOnlySchema && !shouldCreateNewBranch {
		cmd.PrintErrln("'--only-schema' option can only be provided when creating new branches i.e. with '-n, --new' option")
		return
	}
	if shouldCopyOnlySchema && shouldKeepItClean {
		cmd.PrintErrln("'--clean, -c' and '--only-schema' options can not be used simultaneously")
		return
	}

	if shouldCreateNewBranch {
		var err error
		if shouldKeepItClean {
			err = db_int.CreateDb(branchname)
		} else {
			err = lib.NewBranchFromCurrent(branchname, shouldCopyOnlySchema)
		}
		if err != nil {
			if errors.Is(err, errs.ErrDbAlreadyExists) {
				cmd.PrintErrf("Could not create a new database branch with name '%s' because it already exists.\n", branchname)
				return
			}

			cmd.PrintErrln("UNKNOWN ERROR OCCURED!", err)
			return
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
