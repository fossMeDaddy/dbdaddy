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
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wg sync.WaitGroup

	// flags
	titleFlag     string
	dryRunFlag    bool
	emptyInfoFile bool
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
		migrationsDirPath, err := libUtils.GetMigrationsDir(currBranch)
		if err != nil {
			return err
		}

		currentState, err := db_int.GetDbSchema(currBranch)
		if err != nil {
			return err
		}

		migStat, err := migrationsLib.Status(currBranch, currentState)
		if err != nil {
			return err
		}

		var latestMig *types.DbMigration
		prevState := &types.DbSchema{}
		if !migStat.IsInit {
			latestMig = &migStat.Migrations[len(migStat.Migrations)-1]
			state, err := latestMig.ReadState()
			if err != nil {
				return err
			}

			prevState = state
		} else {
			migDirId, err := libUtils.GenerateMigrationId(migrationsDirPath, titleFlag)
			if err != nil {
				return err
			}
			migDirPath := path.Join(migrationsDirPath, migDirId)
			initMig, err := types.NewDbMigration(migDirPath, prevState, "", "", "")
			if err != nil {
				return err
			}

			latestMig = initMig
			// migrations = append(migrations, *latestMig)
		}

		var (
			upSqlScript, downSqlScript string
			upChanges, downChanges     []types.MigAction
		)

		(func() {
			wg.Add(2)
			defer wg.Wait()

			go (func() {
				defer wg.Done()
				upChanges = migrationsLib.DiffDBSchema(currentState, prevState)
				upSqlScript = migrationsLib.GetSQLFromDiffChanges(currentState, prevState, upChanges)
			})()

			go (func() {
				defer wg.Done()
				downChanges = migrationsLib.DiffDBSchema(prevState, currentState)
				downSqlScript = migrationsLib.GetSQLFromDiffChanges(prevState, currentState, downChanges)
			})()
		})()

		if len(upChanges) == 0 {
			cmd.Println("No changes detected.")
			return nil
		}

		if dryRunFlag {
			fmt.Println("-- UP SCRIPT")
			fmt.Println(upSqlScript)

			fmt.Println("-- DOWN SCRIPT")
			fmt.Println(downSqlScript)

			return nil
		}

		migDirId, err := libUtils.GenerateMigrationId(migrationsDirPath, titleFlag)
		if err != nil {
			return err
		}
		migDirPath := path.Join(migrationsDirPath, migDirId)
		mig, err := types.NewDbMigration(
			migDirPath,
			currentState,
			"",
			downSqlScript,
			"",
		)
		if err != nil {
			return err
		}

		if !emptyInfoFile {
			infoFilePath := path.Join(mig.DirPath, constants.MigDirInfoFile)
			libUtils.OpenFileInEditor(infoFilePath)
		}

		if latestMig != nil {
			if err := latestMig.WriteUpQuery(upSqlScript); err != nil {
				return err
			}
		}

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
	cmd.Flags().StringVarP(&titleFlag, "title", "t", "", "add a title for migration file (should not contain any special symbols)")
	cmd.Flags().BoolVar(&emptyInfoFile, "no-info", false, "do not ask for info file input, leave it blank")
	cmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "do not ask for info file input, leave it blank")

	return cmd
}
