package checkoutCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/errs"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

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
	cmd.Flags().BoolVarP(&shouldCreateNewBranch, "new", "n", false, "create a new branch with given branch name, the contents of the current branch will be copied over (tables, columns, column properties, data, etc.)")
	cmd.Flags().BoolVarP(&shouldKeepItClean, "clean", "c", false, "the new branch created by '-n' will be independent of the current branch i.e. nothing will be copied over")
	cmd.Flags().BoolVar(&shouldCopyOnlySchema, "only-schema", false, "create a new branch with given branch name, only the schema of the current branch will be copied over (table structure, relations, etc.)")
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
		var outputFilePath string
		if shouldKeepItClean {
			err = db_int.CreateDb(branchname)
		} else if shouldCopyOnlySchema {
			configFilePath, _ := lib.FindConfigFilePath()

			v := viper.New()
			lib.ReadConfig(v, configFilePath)

			outputFilePath = path.Join(lib.GetDriverDumpDir(configFilePath), fmt.Sprintf("%s__%s", time.Now().Local().Format("2006-01-02_15:04:05"), v.GetString(constants.DbConfigCurrentBranchKey)))
			fmt.Println(outputFilePath)
			if err := db_int.DumpDbOnlySchema(outputFilePath, v); err != nil {
				if err.Error() == errs.PG_DUMP_NOT_FOUND {
					cmd.Println("Hey! we noticed you don't have 'pg_dump', then you also probably won't have 'pg_restore', we use these tools internally to perform dumps & restores... please install these tools in your OS before proceeding.")
				} else {
					cmd.PrintErrln("Unexpected error occured:", err)
				}

				return
			}
			err = db_int.RestoreDbOnlySchema(branchname, viper.GetViper(), outputFilePath)

		} else {
			err = db_int.NewDbFromOriginal(viper.GetString(constants.DbConfigCurrentBranchKey), branchname)
		}
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				cmd.PrintErrf("Could not create a new database branch with name '%s' because it already exists.\n", branchname)
				return
			}

			panic(err)
		}
		if shouldCopyOnlySchema {
			if err := os.RemoveAll(outputFilePath); err != nil {
				panic("Unexpected error occured while removing useless dump files!\n" + err.Error())
			}
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
