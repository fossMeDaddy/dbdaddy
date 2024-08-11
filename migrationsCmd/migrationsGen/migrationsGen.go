package migrationsGenCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files for the current database",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	configDirPath, _ := lib.FindConfigDirPath()
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	migrationsDirPath := path.Join(configDirPath, constants.MigrationDir, currBranch)
	_, err := lib.DirExistsCreate(migrationsDirPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	dirEntries, err := os.ReadDir(migrationsDirPath)
	if err != nil {
		cmd.PrintErrln("unepected error occured while reading migrations dir", migrationsDirPath)
		cmd.PrintErrln(err)
		return
	}

	isInit := false
	if len(dirEntries) == 0 {
		isInit = true
	}

	// var prevState []types.TableSchema

	if isInit {
	} else {
		// get the previous state, parse it & assign it
	}

	// get the current state "getschema" or something

	// upChanges := diff (current_state, previous_state) [GOING UP]
	// downChanges := diff (current_state, previous_state) [GOING DOWN]

	// convert changes to SQL string & write up & down
}

func Init() *cobra.Command {
	// flags

	return cmd
}
