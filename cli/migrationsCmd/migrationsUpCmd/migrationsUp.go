package migrationsUpCmd

import (
	"fmt"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib"
	migrationsLib "github.com/fossmedaddy/dbdaddy/lib/migrationsLib"
	"github.com/fossmedaddy/dbdaddy/middlewares"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "up",
	Short: fmt.Sprintf("Go one migration up by running the '%s' script.", constants.MigDirDownSqlFile),
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.TmpSwitchDB(currBranch, func() error {
		currentState, err := db_int.GetDbSchema()
		if err != nil {
			return err
		}

		migStat, err := migrationsLib.Status(currentState)
		if err != nil {
			return err
		}

		if err := migrationsLib.ApplyMigrationSQL(migStat, true); err != nil {
			return err
		}

		_, currentMigName := path.Split(migStat.ActiveMigration.Up.DirPath)

		cmd.Println("Up migration ran successfully.")
		cmd.Println("Currently at migration state:", currentMigName)
		return nil
	})
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
}

func Init() *cobra.Command {
	// flags

	return cmd
}
