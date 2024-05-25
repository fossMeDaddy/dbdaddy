package restoreCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "restore <dbname>",
	Short: "restores the given database with a database dump",
	Long:  "restoration defaults to restoring the db with the latest dump taken. if no dumps are present with me & i'll require you to provide a path to the backed up dump",
	Run:   cmdRunFn,
}

func getAllDumps() {

}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	configFilePath, _ := lib.FindConfigFilePath()
	dbGroupedDumpFiles, err := lib.GetDbGroupedDumpFiles(configFilePath)
	if err != nil {
		cmd.PrintErrln("Error occured while fetching available dump files!\n", err)
		return
	}
	allDumps := []string{}
	for _, dumpFiles := range dbGroupedDumpFiles {
		allDumps = slices.Concat(allDumps, dumpFiles)
	}

	prompt := promptui.Select{
		Label:             "Available dumps to restore from",
		Items:             allDumps,
		StartInSearchMode: true,
	}

	_, dumpFilePath, err := prompt.Run()
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Printf("Attempting restore on db: '%s' using backup dump file: '%s' ...\n", currBranch, dumpFilePath)
	if err := db_int.RestoreDb(currBranch, viper.GetViper(), dumpFilePath, true); err != nil {
		cmd.PrintErrln("error occured while restoring db\n", err)
		return
	}

	cmd.Println("Success!")
}

func Init() *cobra.Command {
	return cmd
}
