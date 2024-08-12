package migrationsGenCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	migrationsLib "dbdaddy/lib/migrations"
	"dbdaddy/middlewares"
	"dbdaddy/types"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files for the current database",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configDirPath, _ := lib.FindConfigDirPath()
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		migrationsDirPath := path.Join(configDirPath, constants.MigrationDir, currBranch)
		_, err := lib.DirExistsCreate(migrationsDirPath)
		if err != nil {
			return err
		}

		dirEntries, err := os.ReadDir(migrationsDirPath)
		if err != nil {
			return err
		}

		isInit := false
		if len(dirEntries) == 0 {
			isInit = true
		}

		prevState := types.DbSchemaMapping{}
		if !isInit {
			// get the previous state, parse it & assign it
		}

		currentState, err := db_int.GetDbSchemaMapping(currBranch)
		if err != nil {
			return err
		}

		upChanges := migrationsLib.DiffDbSchema(currentState, prevState)
		downChanges := migrationsLib.DiffDbSchema(prevState, currentState)

		fmt.Println(upChanges, downChanges)

		return nil
	})
	if err != nil {
		cmd.PrintErrln("Unexpected error occured!")
		cmd.PrintErrln(err)
		return
	}
}

func Init() *cobra.Command {
	// flags

	return cmd
}
