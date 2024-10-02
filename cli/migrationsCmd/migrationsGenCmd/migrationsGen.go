package migrationsGenCmd

import (
	"path"
	"sync"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	migrationsLib "github.com/fossmedaddy/dbdaddy/lib/migrationsLib"
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wg sync.WaitGroup

	// flags
	titleFlag  string
	dryRunFlag bool
	noInfoFlag bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files for the current database",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.TmpSwitchDB(currBranch, func() error {
		configFilePath, _ := libUtils.FindConfigFilePath()

		migrationsDirPath := libUtils.GetMigrationsDir(path.Dir(configFilePath), currBranch)
		if _, err := libUtils.EnsureDirExists(migrationsDirPath); err != nil {
			return err
		}

		currentState, schemaErr := db_int.GetDbSchema()
		if schemaErr != nil {
			return schemaErr
		}

		latestMig, _, latestMigErr := migrationsLib.GetLatestMigrationOrInit(currentState, titleFlag)
		if latestMigErr != nil {
			return latestMigErr
		}

		prevState, readStateErr := latestMig.ReadState()
		if readStateErr != nil {
			return readStateErr
		}

		var (
			upSqlScript, downSqlScript string
			upChanges, downChanges     []types.MigAction
		)

		var (
			upSqlErr   error
			downSqlErr error
		)
		(func() {
			wg.Add(2)
			defer wg.Wait()

			go (func() {
				defer wg.Done()
				upChanges = migrationsLib.DiffDBSchema(currentState, prevState)
				upSqlScript, upSqlErr = migrationsLib.GetSQLFromDiffChanges(upChanges)
			})()

			go (func() {
				defer wg.Done()
				downChanges = migrationsLib.DiffDBSchema(prevState, currentState)
				downSqlScript, downSqlErr = migrationsLib.GetSQLFromDiffChanges(downChanges)
			})()
		})()
		if downSqlErr != nil {
			cmd.PrintErrln(downSqlErr)
		}
		if upSqlErr != nil {
			cmd.PrintErrln(upSqlErr)
		}
		if downSqlErr != nil || upSqlErr != nil {
			return nil
		}

		if len(upChanges) == 0 {
			cmd.Println("no changes detected.")
			return nil
		}

		if dryRunFlag {
			cmd.Println(constants.DownSqlScriptComment)
			cmd.Println(downSqlScript)

			cmd.Println(constants.UpSqlScriptComment)
			cmd.Println(upSqlScript)

			return nil
		}

		mig, migErr := migrationsLib.GenerateMigration(
			currentState,
			latestMig,
			titleFlag,
			upSqlScript,
			downSqlScript,
			migrationsLib.GetInfoTextFromDiff(upChanges),
		)
		if migErr != nil {
			return migErr
		}

		if !noInfoFlag {
			infoFilePath := path.Join(mig.DirPath, constants.MigDirInfoFile)
			libUtils.OpenFileInEditor(infoFilePath)
		}

		cmd.Println("migration SQL generated successfully.")
		return nil
	})
	if err != nil {
		cmd.PrintErrln("unexpected error occured!")
		cmd.PrintErrln(err)
		return
	}
}

func Init() *cobra.Command {
	// flags
	cmd.Flags().StringVarP(&titleFlag, "title", "t", "", "add a title for migration file (should not contain any special symbols)")
	cmd.Flags().BoolVar(&noInfoFlag, "no-info", false, "do not ask for info file input, leave it blank")
	cmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "print up/down changes console instead of writing them to sql files under migrations dir")

	return cmd
}
