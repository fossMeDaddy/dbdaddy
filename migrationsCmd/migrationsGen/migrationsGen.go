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

		prevState := types.DbSchema{}
		if !isInit {
			// get the previous state, parse it & assign it
		}

		currentState, err := db_int.GetDbSchema(currBranch)
		if err != nil {
			return err
		}

		upChanges := migrationsLib.DiffDbSchema(prevState, currentState)

		upSqlScript := migrationsLib.GetSQLFromDiffChanges(&prevState, &currentState, upChanges)

		fmt.Println(upSqlScript)

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
