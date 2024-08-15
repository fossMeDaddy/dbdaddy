package migrationsGenCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	migrationsLib "dbdaddy/lib/migrations"
	"dbdaddy/libUtils"
	"dbdaddy/middlewares"
	"dbdaddy/types"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	globalWg sync.WaitGroup
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files for the current database",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configDirPath, _ := libUtils.FindConfigDirPath()
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		migrationsDirPath := path.Join(configDirPath, constants.MigrationDir, currBranch)
		_, err := libUtils.DirExistsCreate(migrationsDirPath)
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

		var (
			upSqlScript   string
			downSqlScript string
		)

		(func() {
			globalWg.Add(2)
			defer globalWg.Wait()

			go (func() {
				defer globalWg.Done()
				upChanges := migrationsLib.DiffDbSchema(currentState, prevState)
				upSqlScript = migrationsLib.GetSQLFromDiffChanges(&currentState, &prevState, upChanges)
			})()

			go (func() {
				defer globalWg.Done()
				downChanges := migrationsLib.DiffDbSchema(prevState, currentState)
				downSqlScript = migrationsLib.GetSQLFromDiffChanges(&currentState, &prevState, downChanges)
			})()
		})()

		fmt.Println(upSqlScript, downSqlScript)

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
