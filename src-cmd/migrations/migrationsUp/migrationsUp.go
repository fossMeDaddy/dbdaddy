package migrationsUp

import (
	"dbdaddy/constants"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	migrationsLib "dbdaddy/lib/migrations"
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "up",
	Short: fmt.Sprintf("Go one migration up by running the '%s' script.", constants.MigDirDownSqlFile),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		currentState, err := db_int.GetDbSchema(currBranch)
		if err != nil {
			return err
		}

		migStat, err := migrationsLib.Status(currBranch, currentState)
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
