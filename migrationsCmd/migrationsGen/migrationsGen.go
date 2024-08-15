package migrationsGenCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	migrationsLib "dbdaddy/lib/migrations"
	"dbdaddy/middlewares"
	"dbdaddy/types"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wg sync.WaitGroup
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files for the current database",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		migrations, err := migrationsLib.FetchMigrations(currBranch)
		if err != nil {
			return err
		}

		isInit := false
		if len(migrations) == 0 {
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
			wg.Add(2)
			defer wg.Wait()

			go (func() {
				defer wg.Done()
				upChanges := migrationsLib.DiffDbSchema(currentState, prevState)
				upSqlScript = migrationsLib.GetSQLFromDiffChanges(&currentState, &prevState, upChanges)
			})()

			go (func() {
				defer wg.Done()
				downChanges := migrationsLib.DiffDbSchema(prevState, currentState)
				downSqlScript = migrationsLib.GetSQLFromDiffChanges(&prevState, &currentState, downChanges)
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
