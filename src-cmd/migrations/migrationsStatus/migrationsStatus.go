package migrationsStatusCmd

import (
	"dbdaddy/constants"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	migrationsLib "dbdaddy/lib/migrations"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "status",
	Short: "see migration status, see 'migrations --help' for more info.",
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

		cmd.Println("Database:", currBranch)
		if migStat.ActiveMigration == nil {
			cmd.Println("[WARNING] database schema changed since last migration, please generate migrations via 'migrations generate'")
			return nil
		} else if migStat.ActiveMigration.Up == nil {
			cmd.Println("database schema on latest migration change.")
		}

		cmd.Println()
		cmd.Println("Available migrations:")
		for i, mig := range migStat.Migrations {
			_, migDirName := path.Split(mig.DirPath)
			cmd.Printf("%d. %s ", i+1, migDirName)

			if mig.IsActive {
				cmd.Printf("<--- you are here")
			}

			cmd.Println()
		}

		infoStr, err := migStat.ActiveMigration.GetInfoFile()
		if err != nil {
			return err
		}

		cmd.Println(infoStr)

		return nil
	})
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
}

func Init() *cobra.Command {
	return cmd
}
